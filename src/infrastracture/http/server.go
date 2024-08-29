package http

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/nbisso/storicard-challenge/registry"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(port string) (chan os.Signal, *registry.Registry) {
	reg := registry.NewRegistry()

	go func() {
		r := gin.Default()

		RegisterRoutes(r, reg)

		r.Run(fmt.Sprintf(":%s", port))
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	return c, reg

}
