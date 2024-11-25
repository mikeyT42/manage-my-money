package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const format = "01/02/06"

type BankTransaction struct {
	Id          int32           `csv:"ID" db:"id"`
	Description string          `csv:"Description" db:"description"`
	Type        TransactionType `csv:"Type" db:"type"`
	Amount      Money           `csv:"Amount" db:"amount"`
	Date        Date            `csv:"Date" db:"date"`
	CategoryId  int32           `db:"category"`
}

type Date struct {
	time.Time
}

type Money float32

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
)

// -----------------------------------------------------------------------------
func (a *Money) UnmarshalCSV(csv string) error {
	is_negative := strings.HasPrefix(csv, "($")

	var just_number string
	if is_negative {
		s, _ := strings.CutPrefix(csv, "($")
		s1, _ := strings.CutSuffix(s, ")")
		just_number = "-" + s1
	} else {
		s, _ := strings.CutPrefix(csv, "$")
		just_number = s
	}
	just_number = strings.ReplaceAll(just_number, ",", "")

	f, err := strconv.ParseFloat(just_number, 32)
	if err != nil {
		return err
	}
	*a = Money(float32(f))

	return nil
}

// -----------------------------------------------------------------------------
func (a *Money) MarshalCSV() (string, error) {
	return fmt.Sprintf("%f", float32(*a)), nil
}

// -----------------------------------------------------------------------------
func (a Money) Abs() Money {
	return Money(math.Abs(float64(a)))
}

// -----------------------------------------------------------------------------
func (a Money) ToFloat() float32 {
	return float32(a)
}

// -----------------------------------------------------------------------------
func (tt *TransactionType) UnmarshalCSV(csv string) error {
	var withdrawal_matches = []string{"pos withdrawal", "external withdrawal",
		"withdrawal", "check", "eft credit", "atm withdrawal",
        "insufficient funds charge"}
	var deposit_matches = []string{"deposit", "external deposit", "atm deposit"}

	for _, dm := range deposit_matches {
		if strings.Contains(strings.ToLower(csv), dm) {
			*tt = Deposit
			return nil
		}
	}
	for _, wm := range withdrawal_matches {
		if strings.Contains(strings.ToLower(csv), wm) {
			*tt = Withdrawal
			return nil
		}
	}

	return errors.New("Could not match this transaction type " + csv)
}

// -----------------------------------------------------------------------------
func (tt *TransactionType) MarshalCSV() (string, error) {
	if *tt == Deposit {
		return string(Deposit), nil
	} else {
		return string(Withdrawal), nil
	}
}

// -----------------------------------------------------------------------------
func (d *Date) UnmarshalCSV(csv string) error {
	t, err := time.Parse(format, csv)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// -----------------------------------------------------------------------------
func (d *Date) MarshalCSV() (string, error) {
	return d.Time.Format(format), nil
}

// -----------------------------------------------------------------------------
func (d *Date) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		d.Time = t
	}
	return nil
}

// -----------------------------------------------------------------------------
func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}
