package ws

import (
	"github.com/ghtak/golang.grpc.base/internal/adapter/ginfx"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func NewController() Controller {
	return &controller{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

type Controller interface {
	ginfx.Router
}

type controller struct {
	upgrader websocket.Upgrader
}

func (ct *controller) Register(engine *gin.Engine) error {
	engine.GET("/ws", func(c *gin.Context) {
		conn, err := ct.upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		var (
			mt  int
			msg []byte
		)
		for {
			if mt, msg, err = conn.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %d %s", mt, msg)
			if err = conn.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	})
	return nil
}
