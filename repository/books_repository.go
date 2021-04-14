package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type BooksRepository interface {
	AddBook(book *Book) error
	RemoveBook(id int) error
	GetBook(id int) (*Book, error)
	GetBooks() ([]*Book, error)

}

type sqliteReposotory struct {
	db *sql.DB
}

func (r *sqliteReposotory) AddBook(book *Book) error {
	if _, err := r.db.Exec("INSERT INTO books(name, author) VALUES(?, ?);", book.Name, book.Author); err != nil {
		return err
	}
	return nil
}

func (r *sqliteReposotory) RemoveBook(id int) error {
	if _, err := r.db.Exec("DELETE FROM books WHERE id = ?;", id); err != nil {
		return err
	}
	return nil
}

func (r *sqliteReposotory) GetBook(id int) (*Book, error) {
	row := r.db.QueryRow("SELECT id, name, author FROM books WHERE id = ?;", id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	res := Book{}
	if err := row.Scan(&res.Id, &res.Name, &res.Author); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *sqliteReposotory) 	GetBooks() ([]*Book, error) {
	rows, err := r.db.Query("SELECT id, name, author FROM books;")
	if err != nil {
		return nil, err
	}
	res := make([]*Book, 0)

	for rows.Next() {
		book := Book{}
		if err := rows.Scan(&book.Id, &book.Name, &book.Author); err != nil {
			return nil, err
		}
		res = append(res, &book)
	}

	return res, nil
}

func NewBookRepository(db *sql.DB) BooksRepository {
	repo := &sqliteReposotory{db}
	return repo
}
