package services

import (
	"errors"
	"sync"
	"time"

	"github.com/YoussefMahmod/MoneyTransfering-API/models"
	"github.com/YoussefMahmod/MoneyTransfering-API/store"
	"github.com/google/uuid"
)

type AccountsServiceHandler struct {
	store *store.Datastore
	sync.Mutex
}

func NewAccountsServiceHandler(store *store.Datastore) *AccountsServiceHandler {
	return &AccountsServiceHandler{store: store}
}

func (svc *AccountsServiceHandler) InsertOne(account models.IAccount) {
	svc.store.AccountsByID.Set(account.GetID(), account)
}

func (svc *AccountsServiceHandler) InsertMany(accounts []models.IAccount) {
	var wg sync.WaitGroup
	for idx := range accounts {
		wg.Add(1)

		go func(acc models.IAccount) {
			svc.store.AccountsByID.Set(acc.GetID(), acc)
			wg.Done()
		}(accounts[idx])
	}
	wg.Wait()
}

func (svc *AccountsServiceHandler) GetAll() []models.IAccount {
	accounts := make([]models.IAccount, svc.store.AccountsByID.Count())
	data := svc.store.AccountsByID.GetAll()

	i := 0
	for k := range accounts {
		accounts[i] = data[k].(models.IAccount)
		i++
	}

	return accounts
}

func (svc *AccountsServiceHandler) GetOneByID(id uuid.UUID) (interface{}, bool) {
	accountByID, exists := svc.store.AccountsByID.Get(id)
	return accountByID, exists
}

func (svc *AccountsServiceHandler) PatchOneByID(id uuid.UUID, account models.IAccount) (interface{}, error) {
	_, exists := svc.store.AccountsByID.Get(id)
	if !exists {
		return nil, errors.New("invalid account id")
	}

	var x models.IAccount
	y, _ := svc.store.AccountsByID.Get(id)
	x = y.(models.IAccount)

	x.SetName(account.GetName())
	x.SetBalance(account.GetBalance())
	x.SetUpdatedAt(time.Now())

	return x, nil
}

func (svc *AccountsServiceHandler) DelOneByID(id uuid.UUID) bool {
	_, exist := svc.store.AccountsByID.Del(id)

	if !exist {
		return false
	}

	svc.store.AccountsByID.Del(id)
	return true
}

func (svc *AccountsServiceHandler) GetTxnsBySenderID(id uuid.UUID) (interface{}, bool) {
	txnsBySenderID, exists := svc.store.TransactionsBySender[id]
	return txnsBySenderID, exists
}

func (svc *AccountsServiceHandler) GetTxnsByRecieverID(id uuid.UUID) (interface{}, bool) {
	txnsByRecieverID, exists := svc.store.TransactionsByReciever[id]
	return txnsByRecieverID, exists
}

func (svc *AccountsServiceHandler) CreateTxn(txn models.ITransaction) error {
	acc1, exists := svc.GetOneByID(txn.GetRecieverID())
	if !exists {
		return errors.New("reciever account is not exists")
	}

	acc2, exists := svc.GetOneByID(txn.GetSenderID())
	if !exists {
		return errors.New("sender account is not exists")
	}

	recieverAccount := acc1.(models.IAccount)
	senderAccount := acc2.(models.IAccount)

	smallerIDAcc, largerIDAcc := models.SortAccounts(recieverAccount, senderAccount, false)

	// lock accounts in order based on sorting IDs to avoid deadlocks
	smallerIDAcc.Lock()
	largerIDAcc.Lock()

	// Check if the sender has sufficient balance
	if senderAccount.GetBalance().Cmp(txn.GetAmount()) == -1 {
		return errors.New("insufficient funds")
	}

	// Perform the transfer
	senderAccount.SubBalance(txn.GetAmount())
	recieverAccount.AddBalance(txn.GetAmount())

	// release locks in reverse order
	smallerIDAcc.UnLock()
	largerIDAcc.UnLock()

	// save results (Note that Set could be concurrently but it's thread-safe)
	svc.store.TransactionsByID.Set(txn.GetID(), txn)
	svc.store.TransactionsBySender[senderAccount.GetID()] = append(svc.store.TransactionsBySender[senderAccount.GetID()], txn)
	svc.store.TransactionsByReciever[recieverAccount.GetID()] = append(svc.store.TransactionsByReciever[recieverAccount.GetID()], txn)
	return nil
}
