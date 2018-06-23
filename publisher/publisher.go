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

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/message", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	r.Run(":5000")

	var redis = NewRedisListener().Subscribe("c1", func(message []byte) {
		fmt.Println("message: %v", message)
	})

	<-redis.ready

}
