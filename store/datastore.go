package store

import (
	"github.com/YoussefMahmod/MoneyTransfering-API/models"
	customtypes "github.com/YoussefMahmod/MoneyTransfering-API/utils/custom_types"
	"github.com/google/uuid"
)

type Datastore struct {
	AccountsByID           customtypes.ShardMap
	TransactionsByID       customtypes.ShardMap
	TransactionsBySender   map[uuid.UUID][]models.ITransaction
	TransactionsByReciever map[uuid.UUID][]models.ITransaction
}

func NewDatastore() *Datastore {
	return &Datastore{
		// Index on Accounts(ID)
		AccountsByID: *customtypes.NewShardMap(16),
		// Index on Transactions(ID)
		TransactionsByID: *customtypes.NewShardMap(16),
		// Index on TransactionsBySender(ID)
		TransactionsBySender: make(map[uuid.UUID][]models.ITransaction),
		// Index on TransactionsByReciver(ID)
		TransactionsByReciever: make(map[uuid.UUID][]models.ITransaction),
	}
}
