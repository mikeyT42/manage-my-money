package models

import (
	"time"
)

type Transaction struct {
	Id              int             `db:"id"`
	Type            TransactionType `db:"transaction_type"`
	Amount          Money           `db:"amount"`
	Date            time.Time       `db:"date"`
	BankDescription string          `db:"bank_description"`
	UserDescription string          `db:"user_description"`
	CategoryName    string          `db:"name"`
	CategoryId      int             `db:"category"`
	ColorName       Color           `db:"color"`
}
