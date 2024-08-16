package echo

import (
	"github.com/ghtak/golang.grpc.base/internal/adapter/ginfx"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewController() Controller {
	return controller{}
}

type Controller interface {
	ginfx.Router
}

type controller struct {
}

func (c controller) Register(engine *gin.Engine) error {
	echo := engine.Group("/echo")
	echo.GET("/:echo", c.echo)
	return nil
}

func (c controller) echo(ctx *gin.Context) {
	ctx.String(http.StatusOK, ctx.Param("echo"))
}
