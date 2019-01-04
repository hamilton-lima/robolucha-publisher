package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	redis "gitlab.com/robolucha/robolucha-publisher/redis"
	melody "gopkg.in/olahol/melody.v1"
)

// WatchDetails
type WatchDetails struct {
	MatchID    uint `json:"matchID"`
	LuchadorID uint `json:"luchadorID"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	r := gin.Default()
	m := melody.New()

	var listener = redis.NewRedisListener().SetDebugger(true)
	listener.Connect()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Nothing to see here.",
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	// m.HandleConnect(func(s *melody.Session) {
	// })

	m.HandleDisconnect(func(s *melody.Session) {
		listener.UnSubscribe(s)
	})

	// Watch details information
	m.HandleMessage(func(s *melody.Session, msg []byte) {

		var details WatchDetails
		err := json.Unmarshal(msg, &details)
		if err != nil {
			log.Error("Invalid message content on HandleMessage")
			return
		}

		log.WithFields(log.Fields{
			"details": details,
		}).Info("handleMessage")

		matchStateChannel := fmt.Sprintf("match.%v.state", details.MatchID)
		matchEventChannel := fmt.Sprintf("match.%v.event", details.MatchID)
		luchadorMessageChannel := fmt.Sprintf("luchador.%v.message", details.LuchadorID)

		handler := redis.OnMessageHandler{
			Session: s,
			Handler: func(s *melody.Session, message []byte) {
				s.Write(message)
			}}

		listener.Subscribe(matchStateChannel, handler)
		listener.Subscribe(matchEventChannel, handler)
		listener.Subscribe(luchadorMessageChannel, handler)
	})

	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	r.Run(":5000")

	log.Info("Publisher started on port 5000")
}
