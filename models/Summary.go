package models

type Summary struct {
	CategoryId   int     `db:"id"`
	CategoryName string  `db:"name"`
	Amount       float32 `db:"sum"`
	BudgetAmount float32 `db:"budget_amount"`
	HasSummary   bool    `db:"has_summary"`
}
