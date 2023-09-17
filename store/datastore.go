package store

import (
	customtypes "moneytransfer-api/utils/custom_types"
	"sync"
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
