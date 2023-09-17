package services

import (
	"errors"
	"moneytransfer-api/models"
	"moneytransfer-api/store"
	"time"

	"github.com/google/uuid"
)

type AccountsServiceHandler struct {
	store *store.Datastore
}

func NewAccountsServiceHandler(store *store.Datastore) *AccountsServiceHandler {
	return &AccountsServiceHandler{store: store}
}

func (svc *AccountsServiceHandler) InsertOne(account models.IAccount) {
	svc.store.AccountsByID[account.GetID()] = account
}

func (svc *AccountsServiceHandler) InsertMany(accounts []models.IAccount) {
	for _, account := range accounts {
		svc.store.AccountsByID[account.GetID()] = account
	}
}

func (svc *AccountsServiceHandler) GetAll() []models.IAccount {
	accounts := make([]models.IAccount, len(svc.store.AccountsByID))

	i := 0
	for k := range svc.store.AccountsByID {
		accounts[i] = svc.store.AccountsByID[k]
		i++
	}

	return accounts
}

func (svc *AccountsServiceHandler) GetOneByID(id uuid.UUID) (models.IAccount, bool) {
	accountByID, exists := svc.store.AccountsByID[id]
	return accountByID, exists
}

func (svc *AccountsServiceHandler) PatchOneByID(id uuid.UUID, account models.IAccount) (models.IAccount, error) {
	_, exists := svc.store.AccountsByID[id]
	if !exists {
		return nil, errors.New("invalid account id")
	}

	svc.store.AccountsByID[id].SetName(account.GetName())
	svc.store.AccountsByID[id].SetBalance(account.GetBalance())
	svc.store.AccountsByID[id].SetUpdatedAt(time.Now())

	return svc.store.AccountsByID[id], nil
}

func (svc *AccountsServiceHandler) DelOneByID(id uuid.UUID) bool {
	_, exist := svc.store.AccountsByID[id]

	if !exist {
		return false
	}

	delete(svc.store.AccountsByID, id)
	return true
}
