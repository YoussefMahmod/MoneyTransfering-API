package models

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/YoussefMahmod/MoneyTransfering-API/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const CHUNK_THRESHOLD = 10

type IAccount interface {
	SetDefaults()
	GetID() uuid.UUID
	GetName() string
	SetName(n string)
	GetBalance() decimal.Decimal
	SetBalance(b decimal.Decimal)
	SubBalance(d decimal.Decimal)
	AddBalance(d decimal.Decimal)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
}

type account struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func NewAccount(data []byte) (IAccount, error) {
	var acc account
	acc.SetDefaults()
	err := json.Unmarshal(data, &acc)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func NewListAccounts(data []byte) ([]IAccount, error) {
	var accountsBulk []account
	var wg sync.WaitGroup

	err := json.Unmarshal(data, &accountsBulk)
	if err != nil {
		return nil, err
	}

	accountsBulkLength := len(accountsBulk)
	chunk_length := utils.Max(accountsBulkLength/CHUNK_THRESHOLD, 1)

	res := make([]IAccount, accountsBulkLength)
	for i := 0; i < accountsBulkLength; i += chunk_length {
		wg.Add(1)
		go func(l int, r int) {
			for ; l < r; l++ {
				accountsBulk[l].SetDefaults()
				res[l] = &accountsBulk[l]
			}
			wg.Done()
		}(i, utils.Min(i+chunk_length, accountsBulkLength))
	}

	wg.Wait()

	return res, nil
}

func (a *account) SetDefaults() {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}

	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = a.CreatedAt
	}
}

func (a *account) GetID() uuid.UUID {
	return a.ID
}

func (a *account) GetName() string {
	return a.Name
}

func (a *account) SetName(n string) {
	a.Name = n
}

func (a *account) GetBalance() decimal.Decimal {
	return a.Balance
}

func (a *account) SetBalance(b decimal.Decimal) {
	a.Balance = b
}

func (a *account) SubBalance(b decimal.Decimal) {
	a.Balance.Sub(b)
}

func (a *account) AddBalance(b decimal.Decimal) {
	a.Balance.Add(b)
}

func (a *account) GetCreatedAt() time.Time {
	return a.CreatedAt
}

func (a *account) GetUpdatedAt() time.Time {
	return a.UpdatedAt
}

func (a *account) SetUpdatedAt(t time.Time) {
	a.UpdatedAt = t
}
