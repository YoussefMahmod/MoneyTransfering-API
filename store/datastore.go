package store

import (
	"sync"

	customtypes "github.com/YoussefMahmod/MoneyTransfering-API/utils/custom_types"
)

type Datastore struct {
	AccountsByID customtypes.ShardMap
	Mutex        sync.RWMutex
}

func NewDatastore() *Datastore {
	return &Datastore{
		// Index on Accounts(ID)
		AccountsByID: *customtypes.NewShardMap(16),
		Mutex:        sync.RWMutex{},
	}
}
