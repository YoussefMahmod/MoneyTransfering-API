package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/YoussefMahmod/MoneyTransfering-API/services"
	"github.com/YoussefMahmod/MoneyTransfering-API/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type Server struct {
	Router *gin.Engine
}

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func NewServer() *Server {
	g := gin.Default()

	s := &Server{
		Router: g,
	}

	s.Router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "I am Alive!"})
	})

	dataStore := store.NewDatastore()

	Account{}.router(s, services.NewAccountsServiceHandler(dataStore))
	Transaction{}.router(s, services.NewTransactionsServiceHandler(dataStore))

	return s
}

func (s *Server) Start(port int) { // change to ENV
	s.Router.Run(fmt.Sprintf(":%v", port))
}
