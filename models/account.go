package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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

func NewAccount() account {
	acc := account{}
	acc.SetDefaults()
	return acc
}

func NewListAccounts() []account {
	var acc []account
	return acc
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
