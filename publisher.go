package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	redis "github.com/hamilton-lima/robolucha-publisher/redis"
	session "github.com/hamilton-lima/robolucha-publisher/session"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	m := melody.New()
	var sessionManager = session.NewSessionManager()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "nothing to see here",
		})
	})

	r.GET("/match/:id", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		sessionManager.AddSession(s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		sessionManager.RemoveSession(s)
	})

	//TODO: replace this by ONLY pushing the messages from REDIS to the active sessions
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var matchID = sessionManager.GetIDFromURL(s.Request.URL)
		sessionManager.Broadcast(matchID, msg)
	})

	var listener = redis.NewRedisListener().SetDebugger(true)

	listener.Subscribe("c1", func(message string) {
		fmt.Printf("message: %v \n", message)
	}).Subscribe("c2", func(message string) {
		fmt.Printf("message (c2) : %v \n", message)
	})

	listener.Connect()

	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	r.Run(":5000")

}
