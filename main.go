package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

func main() {
	log.Println("Manage my money...")

	s, err := fs.Sub(static, "static")
	if err != nil {
		panic("Could not open the static directory.")
	}
	t, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("Could not open the templates directory.")
	}

	http.HandleFunc("/manage-my-money/static/styles.css",
        func(w http.ResponseWriter, r *http.Request) {
        http.ServeFileFS(w, r, s, "styles.css")
    })

	http.HandleFunc("/manage-my-money", func(w http.ResponseWriter,
		r *http.Request) {
		log.Println("Loading the main page.")
		tmpl := template.Must(template.ParseFS(t, "index.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc(SummariesRoute,
        func(w http.ResponseWriter, r *http.Request) {
            GetSummaries(w, r, &t)
        })

	http.HandleFunc(UpdateSummaryRoute,
        func(w http.ResponseWriter, r *http.Request) {
            UpdateSummary(w, r, &t)
        })

	http.HandleFunc(CategorizeTransactionPgRoute,
		func(w http.ResponseWriter, r *http.Request) {
			CategorizeTransactionPage(w, r, &t)
		})

	http.HandleFunc(TransactionsUploadRoute,
		func(w http.ResponseWriter, r *http.Request) {
			UploadFile(w, r, &t)
		})

	http.HandleFunc(CustomizeCategoriesPageRoute,
		func(w http.ResponseWriter, r *http.Request) {
			CustomizeCategoriesPage(w, r, &t)
		})

	http.HandleFunc(GetTransactionsRoute,
		func(w http.ResponseWriter, r *http.Request) {
			GetTransactions(w, r, &t)
		})

	http.HandleFunc(EditTransactionsRoute,
		func(w http.ResponseWriter, r *http.Request) {
            EditTransactions(w, r, &t)
		})

	http.HandleFunc(ReadTransactionsButtonRoute,
		func(w http.ResponseWriter, r *http.Request) {
            ReadTransactionsButton(w, r, &t)
		})

	http.HandleFunc(EditTransactionsButtonRoute,
		func(w http.ResponseWriter, r *http.Request) {
            EditTransactionsButton(w, r, &t)
		})

	http.HandleFunc(UpdateTransactionRoute,
		func(w http.ResponseWriter, r *http.Request) {
            UpdateTransaction(w, r, &t)
		})

	http.HandleFunc(EditCategoriesRoute,
		func(w http.ResponseWriter, r *http.Request) {
			EditCategories(w, r, &t)
		})

	http.HandleFunc(DeleteCategoriesRoute,
		func(w http.ResponseWriter, r *http.Request) {
			DeleteCategory(w, r, &t)
		})

    http.HandleFunc(ReadCategoriesRoute,
        func(w http.ResponseWriter, r *http.Request) {
            GetCategories(w, r, &t)
        })

    http.HandleFunc(AddCategoriesRoute,
        func(w http.ResponseWriter, r *http.Request) {
            AddCategory(w, r, &t)
        })

	log.Fatalln(http.ListenAndServe(":8082", nil))
}
