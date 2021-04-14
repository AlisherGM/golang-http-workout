package main

import (
	"book-store/api"
	"book-store/repository"
	"database/sql"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := repository.NewBookRepository(db)

	if err := api.Run(repo); err != nil {
		log.Fatal(err)
	}
}