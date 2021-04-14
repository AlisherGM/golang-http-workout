package api

import (
	"book-store/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type booksApi struct {
	repo repository.BooksRepository
}

func (api *booksApi) getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := func () error {
		books, err := api.repo.GetBooks()
		if err != nil {
			return err
		}
		if err = json.NewEncoder(w).Encode(books); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *booksApi) addBook(w http.ResponseWriter, r *http.Request) {
	err := func () error {
		book := repository.Book{}
		err := json.NewDecoder(r.Body).Decode(&book)

		if err != nil {
			return err
		}

		if err := api.repo.AddBook(&book); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *booksApi) getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		http.NotFound(w, r)
		return
	}
	idVal, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = func () error {
		book, err := api.repo.GetBook(idVal)
		if err != nil {
			return err
		}
		if err = json.NewEncoder(w).Encode(book); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *booksApi) removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		http.NotFound(w, r)
		return
	}
	idVal, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = api.repo.RemoveBook(idVal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func jsonHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Run(repo repository.BooksRepository) error {
	r := mux.NewRouter()
	r.Use(jsonHeaderMiddleware)
	api := booksApi {repo}
	r.HandleFunc("/books/", api.getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", api.getBook).Methods("GET")
	r.HandleFunc("/books/", api.addBook).Methods("POST")
	r.HandleFunc("/books/{id}", api.removeBook).Methods("DELETE")

	return http.ListenAndServe(":8080", r)
}
