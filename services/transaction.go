package services

import (
	"github.com/YoussefMahmod/MoneyTransfering-API/models"
	"github.com/YoussefMahmod/MoneyTransfering-API/store"
	"github.com/google/uuid"
)

type TransactionsServiceHandler struct {
	store *store.Datastore
}

func NewTransactionsServiceHandler(store *store.Datastore) *TransactionsServiceHandler {
	return &TransactionsServiceHandler{store: store}
}

func (svc *TransactionsServiceHandler) GetAll() []models.ITransaction {
	transactions := make([]models.ITransaction, svc.store.TransactionsByID.Count())
	data := svc.store.TransactionsByID.GetAll()

	i := 0
	for k := range transactions {
		transactions[i] = data[k].(models.ITransaction)
		i++
	}

	return transactions
}

func (svc *TransactionsServiceHandler) GetOneByID(id uuid.UUID) (interface{}, bool) {
	transactionByID, exists := svc.store.TransactionsByID.Get(id)
	return transactionByID, exists
}

func (svc *TransactionsServiceHandler) DelOneByID(id uuid.UUID) bool {
	_, exist := svc.store.TransactionsByID.Del(id)

	if !exist {
		return false
	}

	svc.store.TransactionsByID.Del(id)
	return true
}
