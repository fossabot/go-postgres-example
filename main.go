package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
)

type Book struct {
	Title       string
	Author      string
	Description string
}

func SetupDB() *sql.DB {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", connectionString)
	PanicIf(err)
	return db
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	m := martini.Classic()
	m.Map(SetupDB())
	m.Use(render.Renderer(
		render.Options{
			Layout: "layout",
		}))

	m.Get("/", func(ren render.Render, r *http.Request, db *sql.DB) {
		searchTerm := "%" + r.URL.Query().Get("search") + "%"
		rows, err := db.Query("SELECT title, author, description FROM books WHERE title ILIKE $1 OR author ILIKE $1 OR description ILIKE $1", searchTerm)
		PanicIf(err)
		defer rows.Close()

		books := []Book{}

		for rows.Next() {
			PanicIf(rows.Err())
			b := Book{}
			err := rows.Scan(&b.Title, &b.Author, &b.Description)
			PanicIf(err)
			books = append(books, b)
		}

		ren.HTML(200, "books", books)
	})

	m.Run()
}
