package publisher

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	m := melody.New()
	var sessionManager = NewSessionManager()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
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

	var redis = NewRedisListener().SetDebugger(true)

	redis.Subscribe("c1", func(message string) {
		fmt.Printf("message: %v \n", message)
	}).Subscribe("c2", func(message string) {
		fmt.Printf("message (c2) : %v \n", message)
	})

	redis.Connect()

	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	r.Run(":5000")

}
