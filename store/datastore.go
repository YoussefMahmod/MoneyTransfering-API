package store

import (
	"moneytransfer-api/models"
	"sync"

	"github.com/google/uuid"
)

type Datastore struct {
	AccountsByID map[uuid.UUID]models.IAccount
	Mutex        sync.RWMutex
}

func NewDatastore() *Datastore {
	return &Datastore{
		// Index on Accounts(ID)
		AccountsByID: make(map[uuid.UUID]models.IAccount),
		Mutex:        sync.RWMutex{},
	}
}
