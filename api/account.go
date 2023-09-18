package api

import (
	"fmt"
	"io"
	"moneytransfer-api/models"
	"moneytransfer-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Account struct {
	server         *Server
	serviceHandler *services.AccountsServiceHandler
}

func (u Account) router(server *Server, svcHandler *services.AccountsServiceHandler) {
	u.server = server
	u.serviceHandler = svcHandler

	serverGroup := server.Router.Group("/accounts")
	serverGroup.GET("/", u.getAccountsList)
	serverGroup.GET("/:id", u.getAccount)
	serverGroup.POST("/", u.createAccount)
	serverGroup.POST("/bulk", u.createAccountsInBulk)
	serverGroup.PATCH("/:id", u.patchAccount)
	serverGroup.DELETE("/:id", u.delAccount)
}

func (a *Account) getAccountsList(c *gin.Context) {
	accounts := a.serviceHandler.GetAll()

	c.JSON(http.StatusOK, accounts)
}

func (a *Account) getAccount(c *gin.Context) {
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

func (a *Account) createAccount(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := models.NewAccount(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.serviceHandler.InsertOne(account)

	c.JSON(http.StatusCreated, account)
}

func (a *Account) createAccountsInBulk(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accounts, err := models.NewListAccounts(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.serviceHandler.InsertMany(accounts)

	Logger.Info(fmt.Sprintf("INFO - %v accounts are ingested and ready to transfer", len(accounts)))

	c.JSON(http.StatusCreated, accounts)
}

func (a *Account) patchAccount(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := models.NewAccount(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAccount, err := a.serviceHandler.PatchOneByID(uuid, account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newAccount)
}

func (a *Account) delAccount(c *gin.Context) {
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
