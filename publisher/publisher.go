package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hamilton-lima/robolucha-services/redislistener"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/message", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	var redis = redislistener.NewRedisListener().SetDebugger(true)

	redis.Subscribe("c1", func(message string) {
		fmt.Printf("message: %v \n", message)
	}).Subscribe("c2", func(message string) {
		fmt.Printf("message (c2) : %v \n", message)
	})

	redis.Connect()

	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	r.Run(":5000")

}
