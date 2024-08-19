package ws

import (
	"github.com/ghtak/golang.grpc.base/internal/adapter/fiberfx"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Controller interface {
	fiberfx.Router
}

type wsController struct {
}

func (c wsController) Register(router fiber.Router) error {
	router.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// "/ws/:id" -> log.Println(c.Params("id"))
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %d %s", mt, msg)
			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))
	return nil
}

func NewController() Controller {
	return wsController{}
}
