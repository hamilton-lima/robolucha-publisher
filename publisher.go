package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
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

func setLogLevel(name string) {
	level := os.Getenv(name)
	if level == "" {
		level = "info"
	}

	logLevel, err := log.ParseLevel(level)

	if err != nil {
		log.WithFields(log.Fields{
			"level": level,
		}).Warning("Invalid log level, default 'info'")
		log.SetLevel(log.InfoLevel)
	} else {
		log.WithFields(log.Fields{
			"level": level,
		}).Info("Set Log level")
		log.SetLevel(logLevel)
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	setLogLevel("PUBLISHER_LOG_LEVEL")
	log.Info("Starting Robolucha publisher.")

	r := gin.Default()
	m := melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool {
		log.Debug("CheckOrigin function")
		return true
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

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

	m.HandleConnect(func(s *melody.Session) {
		log.Debug("Handle Connect")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		listener.UnSubscribeAll(s)
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

		listener.UnSubscribeAll(handler.Session)
		listener.Subscribe(matchStateChannel, &handler)
		listener.Subscribe(matchEventChannel, &handler)
		listener.Subscribe(luchadorMessageChannel, &handler)
	})

	r.Run(":5000")

	log.Info("Publisher started on port 5000")
}
