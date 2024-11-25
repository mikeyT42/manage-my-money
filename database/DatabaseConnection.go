package database

import (
	"log"
	"math"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	m "github.com/mikeyT42/manage-my-money/models"
)

var logE = log.New(os.Stderr, "database: ", log.Lshortfile|log.LstdFlags)

// -----------------------------------------------------------------------------
func ConnectToDB() *sqlx.DB {
	dbConfig, err := os.ReadFile("/etc/manage-my-money/db.conf")
	if err != nil {
		logE.Fatalln(err)
	}

	ps, err := sqlx.Connect("postgres", string(dbConfig))
	if err != nil {
		logE.Fatalln(err)
	}

	if err := ps.Ping(); err != nil {
		logE.Fatalln(err)
	}

	return ps
}

// -----------------------------------------------------------------------------
func AddBankStatement(bs m.BankStatement) error {
	db := ConnectToDB()
	defer db.Close()

	i := `insert into transactions (bank_id, "transaction_type", amount, date,
                bank_description, category)
                values (:id, :type, :amount, :date, :description, :category)
                on conflict (bank_id) do nothing`
	stmt, err := db.PrepareNamed(i)
	if err != nil {
		logE.Fatal("Could not prepare the statement for inserting into the " +
			"transactions table.")
	}
	defer stmt.Close()

	for _, t := range bs.Transactions {
		_, err := stmt.Exec(t)
		if err != nil {
			logE.Print("Could not add this transaction to DB:")
			logE.Print(t)
			logE.Print(err)
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
func AddCategory(c m.Category) error {
	db := ConnectToDB()
	defer db.Close()

	i := `insert into categories (description, "name", color)
            values (:description, :name, :color)`
	stmt, err := db.PrepareNamed(i)
	if err != nil {
		logE.Fatal("Could not prepare the statement for inserting into the " +
			"categories table.")
	}
	defer stmt.Close()

	_, err = stmt.Exec(c)
	if err != nil {
		logE.Print("Could not add this category to DB:")
		logE.Print(c)
		logE.Print(err)
	}

	return nil
}

// -----------------------------------------------------------------------------
func DeleteCategory(c m.Category) error {
	db := ConnectToDB()
	defer db.Close()

	d := `delete from categories where "name" = :name`
	stmt, err := db.PrepareNamed(d)
	if err != nil {
		logE.Fatal("Could not prepare the statement for deleting from the " +
			"categories table.")
	}
	defer stmt.Close()

	_, err = stmt.Exec(c)
	if err != nil {
		logE.Print("Could not delete this category from DB:")
		logE.Print(c)
		logE.Print(err)
	}

	return nil
}

// -----------------------------------------------------------------------------
func GetCategories() []m.Category {
	ps := ConnectToDB()
	defer ps.Close()

	q := `select c.id, description, "name", cc."color", c."color" as color_id,
            is_user_defined 
            from categories c 
            join category_colors cc on c."color"=cc.id`
	rows, _ := ps.Queryx(q)
	defer rows.Close()

	var categories []m.Category
	for rows.Next() {
		category := m.Category{}
		err := rows.StructScan(&category)
		if err != nil {
			logE.Print(err)
		}
		categories = append(categories, category)
	}

	return categories
}

// -----------------------------------------------------------------------------
func GetTransactions(page, size int) ([]m.Transaction, int, error) {
	db := ConnectToDB()
	defer db.Close()

	offset := (page - 1) * size
	q := `select t.id, transaction_type, amount, date,
            bank_description, user_description, c."name", cc."color", t.category
            from transactions t 
            join categories c on t.category = c.id
            join category_colors cc on c."color" = cc.id
            order by date desc, bank_id desc
            limit $1
            offset $2`
	rows, err := db.Queryx(q, size, offset)
	if err != nil {
		logE.Println("Could not execute transaction query.")
		logE.Println(err)
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []m.Transaction
	for rows.Next() {
		var transaction m.Transaction
		err := rows.StructScan(&transaction)
		if err != nil {
			logE.Print(err)
		}
		transactions = append(transactions, transaction)
	}

	var total int
	err = db.QueryRowx("select count(*) from transactions").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// -----------------------------------------------------------------------------
func UpdateTransaction(t m.Transaction) error {
	db := ConnectToDB()
	defer db.Close()

	u := `update transactions
            set user_description = :user_description,
                category = :category,
                date = :date
            where id = :id`
	stmt, err := db.PrepareNamed(u)
	if err != nil {
		logE.Print("Could not prepare the statement for updating the " +
			"transaction")
		logE.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(t)
	if err != nil {
		logE.Print("Could not update transaction:")
		logE.Print(t)
		logE.Print(err)
	}

	return nil
}

// -----------------------------------------------------------------------------
func GetColors() []m.CategoryColor {
	db := ConnectToDB()
	defer db.Close()

	q := `select "color", id
            from category_colors`
	rows, _ := db.Queryx(q)
	defer rows.Close()

	var colors []m.CategoryColor
	for rows.Next() {
		var color m.CategoryColor
		err := rows.StructScan(&color)
		if err != nil {
			logE.Print(err)
		}
		colors = append(colors, color)
	}

	return colors
}

// -----------------------------------------------------------------------------
func GetSummaries(month int32, year int32) []m.Summary {
	db := ConnectToDB()
	defer db.Close()

	q := `select 
            sum(t.amount) as sum,
            case when s.budget_amount is null
                then 0
                else s.budget_amount
                end as budget_amount,
            c.name,
            c.id,
            case when s.summary_count > 0
                then true
                else false
                end as has_summary
        from
            transactions t
        join
            categories c on t.category = c.id
        left join (
            select
                category,
                budget_amount,
                count(id) as summary_count
            from
                summaries
            where
            "month" = :month
            group by
                category, budget_amount
        ) s on c.id = s.category
        where
        extract(month from t."date") = :month
        and extract(year from t."date") = :year
        group by
            c.name, c.id, s.summary_count, s.budget_amount
        order by
            c.name`
	stmt, err := db.PrepareNamed(q)
	if err != nil {
		logE.Print(err)
	}
	defer stmt.Close()
	rows, err := stmt.Queryx(map[string]any{
		"month": month,
		"year":  year,
	})
	if err != nil {
		logE.Print(err)
	}
	defer rows.Close()

	var summaries []m.Summary
	for rows.Next() {
		var summary m.Summary
		err := rows.StructScan(&summary)
		if err != nil {
			logE.Print(err)
		}
		summary.Amount = float32(math.Abs(float64(summary.Amount)))
		summaries = append(summaries, summary)
	}

	return summaries
}

// -----------------------------------------------------------------------------
func UpdateSummary(category_id int64, has_summary bool, amount float32,
	month int32, year int32) {
	db := ConnectToDB()
	var sql string
	if has_summary {
		sql = GetUpdateSummarySQL()
	} else {
		sql = GetInsertSummarySQL()
	}
	log.Print(sql)

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		logE.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(map[string]any{
		"category_id": category_id,
		"amount":      amount,
		"month":       month,
		"year":        year,
	})
	if err != nil {
		logE.Print(err)
	}
}

// -----------------------------------------------------------------------------
func GetInsertSummarySQL() string {
	return `
    insert into summaries
        (category, budget_amount, "month", "year")
    values
        (:category_id, :amount, :month, :year)`
}

// -----------------------------------------------------------------------------
func GetUpdateSummarySQL() string {
	return `
        update summaries
        set budget_amount = :amount
        where category = :category_id and "month" = :month and "year" = :year`
}
