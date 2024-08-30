package domain

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	ID       string    `json:"id" csv:"id" db:"id"`
	UserID   string    `json:"user_id" csv:"user_id" db:"user_id"`
	Amount   float64   `json:"amount" csv:"amount" db:"amount"`
	Datetime time.Time `json:"datetime" csv:"datetime" db:"datetime"`
}

type TransactionResult struct {
	Balance      *float64 `json:"balance" db:"balance" default:"0"`
	TotalDebits  *float64 `json:"total_debits" db:"total_debits" default:"0"`
	TotalCredits *float64 `json:"total_credits" db:"total_credits" default:"0"`
}

type TransactionFilter struct {
	UserID string     `json:"user_id" csv:"user_id" db:"user_id" validate:"required"`
	From   *time.Time `json:"from_date" csv:"datetime" db:"fom_datetime" validate:"datetime" form:"from"`
	To     *time.Time `json:"to_date" csv:"datetime" db:"to_datetime" validate:"datetime" form:"to"`
}

func (t *Transaction) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

func NewTransactionEventFromJson(jsonString string) (*Transaction, error) {
	n := &Transaction{}

	err := json.Unmarshal([]byte(jsonString), n)

	return n, err
}
