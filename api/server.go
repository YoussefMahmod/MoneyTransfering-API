package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type Server struct{
	router *gin.Engine
}

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func NewServer() *Server{
	g := gin.Default()
	
	return &Server{
		router: g,
	}
}

func (s *Server) Start(port int) {// change to ENV
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "I am Alive!"})
	})
	
	s.router.Run(fmt.Sprintf(":%v", port))
}