package api

import (
	"net/http"

	"github.com/YoussefMahmod/MoneyTransfering-API/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Transaction struct {
	server         *Server
	serviceHandler *services.TransactionsServiceHandler
}

func (u Transaction) router(server *Server, svcHandler *services.TransactionsServiceHandler) {
	u.server = server
	u.serviceHandler = svcHandler

	serverGroup := server.Router.Group("/api/v1/transactions")
	serverGroup.GET("/", u.getAllTransactions)
	serverGroup.GET("/:id", u.getTransactionByID)
	serverGroup.DELETE("/:id", u.deleteTransactionByID)
}

func (a *Transaction) getAllTransactions(c *gin.Context) {
	accounts := a.serviceHandler.GetAll()

	c.JSON(http.StatusOK, accounts)
}

func (a *Transaction) getTransactionByID(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, exists := a.serviceHandler.GetOneByID(uuid)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Element Not Found!"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (a *Transaction) deleteTransactionByID(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists := a.serviceHandler.DelOneByID(uuid)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Element Not Found!"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
