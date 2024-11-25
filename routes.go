package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	db "github.com/mikeyT42/manage-my-money/database"
	m "github.com/mikeyT42/manage-my-money/models"
	p "github.com/mikeyT42/manage-my-money/pagination"
)

var logE = log.New(os.Stderr, "routes: ", log.Lshortfile|log.LstdFlags)

// -----------------------------------------------------------------------------
//
//	Main Screen
//
// -----------------------------------------------------------------------------
const SummariesRoute = "/manage-my-money/get-summaries"
const UpdateSummaryRoute = "/manage-my-money/update-summary"

// -----------------------------------------------------------------------------
//
//	Customize Categories Screen
//
// -----------------------------------------------------------------------------
const CustomizeCategoriesPageRoute = "/manage-my-money/customize-categories"
const EditCategoriesRoute = "/manage-my-money/edit-mode-categories"
const ReadCategoriesRoute = "/manage-my-money/read-mode-categories"
const AddCategoriesRoute = "/manage-my-money/add-category"
const DeleteCategoriesRoute = "/manage-my-money/delete-category"

// -----------------------------------------------------------------------------
//
//	Customize Transactions Screen
//
// -----------------------------------------------------------------------------
const CategorizeTransactionPgRoute = "/manage-my-money/categorize-transactions"
const TransactionsUploadRoute = "/manage-my-money/upload"
const GetTransactionsRoute = "/manage-my-money/read-mode-transactions"
const EditTransactionsRoute = "/manage-my-money/edit-mode-transactions"
const ReadTransactionsButtonRoute = "/manage-my-money/read-mode-button"
const EditTransactionsButtonRoute = "/manage-my-money/edit-mode-button"
const UpdateTransactionRoute = "/manage-my-money/update-transaction"

// -----------------------------------------------------------------------------
func CategorizeTransactionPage(w http.ResponseWriter, r *http.Request,
	t *fs.FS) {
	log.Println("Loading the Categorize Transactions page.")
	tmpl := template.Must(template.ParseFS(*t, "categorize-transactions.html"))
	tmpl.Execute(w, nil)
}

// -----------------------------------------------------------------------------
func EditCategories(w http.ResponseWriter, r *http.Request,
	t *fs.FS) {
	log.Println("Loading the Edit Categories fragment.")
	categories := db.GetCategories()
	colors := db.GetColors()
	data := map[string]interface{}{
		"Categories": categories,
		"Colors":     colors,
	}
	tmpl := template.Must(template.ParseFS(*t,
		"fragments/edit-mode-categories.html"))
	tmpl.Execute(w, data)
}

// -----------------------------------------------------------------------------
func GetCategories(w http.ResponseWriter, r *http.Request,
	t *fs.FS) {
	log.Println("Loading the Read Custom Categories fragment.")
	categories := db.GetCategories()
	tmpl := template.Must(template.ParseFS(*t,
		"fragments/read-mode-categories.html"))
	tmpl.Execute(w, categories)
}

// -----------------------------------------------------------------------------
func AddCategory(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Print("Adding a category.")
	r.ParseForm()

	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
	}

	category := m.Category{
		Description: r.FormValue("category-description"),
		Name:        r.FormValue("category-name"),
		Color:       r.FormValue("category-color"),
	}
	db.AddCategory(category)

	categories := db.GetCategories()
	colors := db.GetColors()
	data := map[string]interface{}{
		"Categories": categories,
		"Colors":     colors,
	}
	tmpl := template.Must(template.ParseFS(*t,
		"fragments/edit-mode-categories.html"))
	tmpl.Execute(w, data)
}

// -----------------------------------------------------------------------------
func DeleteCategory(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Print("Deleting a category.")
	r.ParseForm()

	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
	}

	category := m.Category{
		Description: r.FormValue("category-description"),
		Name:        r.FormValue("category-name"),
		Color:       r.FormValue("category-color"),
	}
	log.Print(category)

	db.DeleteCategory(category)

	categories := db.GetCategories()
	colors := db.GetColors()
	data := map[string]interface{}{
		"Categories": categories,
		"Colors":     colors,
	}
	tmpl := template.Must(template.ParseFS(*t,
		"fragments/edit-mode-categories.html"))
	tmpl.Execute(w, data)
}

// -----------------------------------------------------------------------------
func CustomizeCategoriesPage(w http.ResponseWriter, r *http.Request,
	t *fs.FS) {
	log.Println("Loading the Customize Categories page.")
	categories := db.GetCategories()
	tmpl := template.Must(template.ParseFS(*t, "customize-categories.html"))
	tmpl.Execute(w, categories)
}

// -----------------------------------------------------------------------------
func UploadFile(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Println("Adding bank transactions.")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		logE.Println("Had trouble parsing the form for the file.")
		logE.Println(err)
	}

	file, fileHeader, err := r.FormFile("statement-upload")
	if err != nil {
		logE.Println("Error reading the file.")
		logE.Println(err)
	}
	defer file.Close()

	b := make([]byte, fileHeader.Size)
	_, err = file.Read(b)
	if err != nil {
		logE.Println(err)
	}

	var bs m.BankStatement
	if err := m.Unmarshal(b, &bs); err != nil {
		logE.Println("Could not unmarshal.")
	}

	bs.AutoCategorize()

	db.AddBankStatement(bs)
	w.Header().Add("HX-Trigger", "bankstatement-upload")
	tmpl := template.Must(template.ParseFS(*t,
		"fragments/bankstatement-upload-input.html"))
	tmpl.Execute(w, nil)
}

// -----------------------------------------------------------------------------
func GetTransactions(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Print("Reading transactions.")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
		logE.Printf("The page query argument could not be grabbed. Using a "+
			"page value of %d.", page)
	}
	if page < 1 {
		page = 1
	}
	size := 25

	transactions, total, err := db.GetTransactions(page, size)
	if err != nil {
		logE.Println("Get not get transactions from DB.")
	}

	w.Header().Add("HX-Trigger", "edit-mode-button")
	tmpl := template.Must(template.ParseFS(*t,
		"fragments/read-mode-transactions-table.html"))
	p := p.PaginationState{
		Previous:         page - 1,
		Next:             page + 1,
		DisabledPrevious: (page - 1) < 1,
		DisabledNext:     (page * size) >= total,
	}

	data := map[string]interface{}{
		"Transactions":    transactions,
		"PaginationState": p,
	}

	log.Println(p)
	tmpl.Execute(w, data)
}

// -----------------------------------------------------------------------------
func ReadTransactionsButton(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Print("Sending the \"Read Mode\" button.")
	tmpl := template.Must(template.ParseFS(*t,
		"fragments/read-mode-button-transactions.html"))
	tmpl.Execute(w, nil)
}

// -----------------------------------------------------------------------------
func EditTransactionsButton(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Print("Sending the \"Edit Mode\" button.")
	tmpl := template.Must(template.ParseFS(*t,
		"fragments/edit-mode-button-transactions.html"))
	tmpl.Execute(w, nil)
}

// -----------------------------------------------------------------------------
func EditTransactions(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Print("Editing transactions.")
	r.ParseForm()

	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
		logE.Printf("The page query argument could not be grabbed. Using a "+
			"page value of %d.", page)
	}
	if page < 1 {
		page = 1
	}
	size := 25

	transactions, total, err := db.GetTransactions(page, size)
	if err != nil {
		logE.Println("Could not get the transactions from the DB.")
	}
	p := p.PaginationState{
		Previous:         page - 1,
		Next:             page + 1,
		DisabledPrevious: (page - 1) < 1,
		DisabledNext:     (page * size) >= total,
	}

	categories := db.GetCategories()

	data := map[string]interface{}{
		"Transactions":    transactions,
		"Categories":      categories,
		"PaginationState": p,
	}

	w.Header().Add("HX-Trigger", "read-mode-button")
	tmpl := template.Must(template.ParseFS(*t,
		"fragments/edit-mode-transactions-table.html"))
	tmpl.Execute(w, data)
}

// -----------------------------------------------------------------------------
func UpdateTransaction(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Print("Updating Transaction.")
	r.ParseForm()
	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
	}

	id, err := strconv.ParseInt(r.FormValue("transaction-id"), 10, 64)
	if err != nil {
		logE.Print(err)
	}
	c_id, err := strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		logE.Print(err)
	}
	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		logE.Print(err)
	}

	transaction := m.Transaction{
		Id:              int(id),
		UserDescription: r.FormValue("user-description"),
		CategoryId:      int(c_id),
		Date:            date,
	}

	db.UpdateTransaction(transaction)
}

// -----------------------------------------------------------------------------
func GetSummaries(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Print("Getting summaries.")

	r.ParseForm()
	month, err := strconv.ParseInt(r.FormValue("month"), 10, 32)
	if err != nil {
		logE.Print(err)
	}
	year, err := strconv.ParseInt(r.FormValue("year"), 10, 32)
	if err != nil {
		logE.Print(err)
	}

	summaries := db.GetSummaries(int32(month), int32(year))

	tmpl := template.Must(template.ParseFS(*t, "fragments/summaries.html"))
	tmpl.Execute(w, summaries)
}

// -----------------------------------------------------------------------------
func UpdateSummary(w http.ResponseWriter, r *http.Request, t *fs.FS) {
	log.Print("Updating summary.")

	r.ParseForm()
	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
	}
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 32)
	if err != nil {
		// TODO: Validation somehow. Maybe through the UI using htmx.
		// I remember seeing something in the docs.
		logE.Print(err)
	}
	category, err := strconv.ParseInt(r.FormValue("category-id"), 10, 64)
	if err != nil {
		logE.Print(err)
	}
	has_summary, err := strconv.ParseBool(r.FormValue("has-summary"))
	if err != nil {
		logE.Print(err)
	}
	month, err := strconv.ParseInt(r.FormValue("month"), 10, 32)
	if err != nil {
		logE.Print(err)
	}
	year, err := strconv.ParseInt(r.FormValue("year"), 10, 32)
	if err != nil {
		logE.Print(err)
	}

	db.UpdateSummary(category, has_summary, float32(amount), int32(month),
		int32(year))
	// TODO: Send an event trigger to show a check mark on success.
}
