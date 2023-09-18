package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ITransaction interface {
	SetDefaults()
	GetID() uuid.UUID
	GetSenderID() uuid.UUID
	SetSenderID(id uuid.UUID)
	GetRecieverID() uuid.UUID
	GetAmount() decimal.Decimal
	SetAmount(b decimal.Decimal)
	SubAmount(d decimal.Decimal)
	AddAmount(d decimal.Decimal)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
}

type transaction struct {
	ID         uuid.UUID       `json:"id"`
	SenderID   uuid.UUID       `json:"sender_id"`
	RecieverID uuid.UUID       `json:"reciever_id"`
	Amount     decimal.Decimal `json:"amount"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func NewTransaction(data []byte) (ITransaction, error) {
	var txn transaction
	txn.SetDefaults()
	err := json.Unmarshal(data, &txn)
	if err != nil {
		return nil, err
	}

	return &txn, nil
}

func (txn *transaction) SetDefaults() {
	if txn.ID == uuid.Nil {
		txn.ID = uuid.New()
	}

	if txn.CreatedAt.IsZero() {
		txn.CreatedAt = time.Now()
	}

	if txn.UpdatedAt.IsZero() {
		txn.UpdatedAt = txn.CreatedAt
	}
}

func (txn *transaction) GetID() uuid.UUID {
	return txn.ID
}

func (txn *transaction) GetSenderID() uuid.UUID {
	return txn.SenderID
}

func (txn *transaction) SetSenderID(id uuid.UUID) {
	txn.SenderID = id
}

func (txn *transaction) GetRecieverID() uuid.UUID {
	return txn.RecieverID
}

func (txn *transaction) GetAmount() decimal.Decimal {
	return txn.Amount
}

func (txn *transaction) SetAmount(b decimal.Decimal) {
	txn.Amount = b
}

func (txn *transaction) SubAmount(b decimal.Decimal) {
	txn.Amount.Sub(b)
}

func (txn *transaction) AddAmount(b decimal.Decimal) {
	txn.Amount.Add(b)
}

func (txn *transaction) GetCreatedAt() time.Time {
	return txn.CreatedAt
}

func (txn *transaction) GetUpdatedAt() time.Time {
	return txn.UpdatedAt
}

func (txn *transaction) SetUpdatedAt(t time.Time) {
	txn.UpdatedAt = t
}

// func (txn *transaction) SortedSenderRecieverIDs
