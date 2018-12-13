package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Структура ответа от сервера
type Response struct {
	Message interface{}
	Error   string
}

type Book struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type BookStore struct {
	books []Book
}

///////////////////////////////////////////////////////////////////

var bookStore = BookStore{
	books: make([]Book, 0),
}

///////////////////////////////////////////////////////////////////

func main() {
	fmt.Println("Listening port on :3000")

	handler := http.NewServeMux()
	handler.HandleFunc("/hello/", Logger(handlerHello))
	handler.HandleFunc("/book/", Logger(handlerBook))
	handler.HandleFunc("/books/", Logger(handlerBooks))

	s := http.Server{
		Addr:           ":3000",          // Адресс сервера
		Handler:        handler,          // Диспетчер маршрутов
		ReadTimeout:    10 * time.Second, // Время чтения запроса
		WriteTimeout:   10 * time.Second, // Время записи ответа
		IdleTimeout:    10 * time.Second, // Время ожидания следующего запроса
		MaxHeaderBytes: 1 << 20,          // Максимальный размер http header в байтах (1 * 2 ^ 20 = 128 kByte)
	}

	log.Fatal(s.ListenAndServe()) // Если есть ошибки то отобразить их
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	name := strings.Replace(r.URL.Path, "/hello/", "", 1)
	response := Response{
		Message: fmt.Sprintf("Hello %s! Glad to see you again.", name),
	}

	respJson, _ := json.Marshal(response)
	_, _ = w.Write(respJson)
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server [net/http] method [%s] connection from [%v]", r.Method, r.RemoteAddr)

		next.ServeHTTP(w, r)
	}
}

///////////////////////////////////////////////////////////////////

func handlerBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet:
		handleGetBook(w, r)

	case r.Method == http.MethodPost:
		handleAddBook(w, r)

	case r.Method == http.MethodPut:
		handleUpdateBook(w, r)

	case r.Method == http.MethodDelete:
		handleDeleteBook(w, r)
	}
}

func handlerBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if r.Method == http.MethodGet {
		handleGetBook(w, r)
	}

	response := Response{
		Message: bookStore.GetBooks(),
	}
	respJson, _ := json.Marshal(response)

	_, _ = w.Write(respJson)
}

///////////////////////////////////////////////////////////////////

func handleGetBook(w http.ResponseWriter, r *http.Request) {
	var response Response
	var status int

	bookId := strings.Replace(r.URL.Path, "/book/", "", 1)
	book := bookStore.FindBookById(bookId)

	if book == nil {
		status = http.StatusNotFound
		response.Error = fmt.Sprintf("Book with Id:%s not found!", bookId)
	} else {
		status = http.StatusOK
		response.Message = book
	}

	w.WriteHeader(status)
	respJson, _ := json.Marshal(response)
	_, _ = w.Write(respJson)
}

func handleAddBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var response Response

	// Считываем книгу
	decoder := json.NewDecoder(r.Body) // Декодируем тело запроса

	// Ошибка в декодере
	err := decoder.Decode(&book) // Передаем адресс на книгу
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = err.Error()

		respJson, _ := json.Marshal(response)
		_, _ = w.Write(respJson)

		return
	}

	// Ошибка при добавлении книги
	err = bookStore.AddBook(book)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		response.Error = err.Error()

		respJson, _ := json.Marshal(response)
		_, _ = w.Write(respJson)

		return
	}

	w.WriteHeader(http.StatusOK)

	response = Response{
		Message: fmt.Sprintf("Book with id:%s updated is SUCCESS!", book.Id),
	}
	respJson, _ := json.Marshal(response)

	_, _ = w.Write(respJson)
}

func handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var response Response

	// Считываем книгу
	decoder := json.NewDecoder(r.Body) // Декодируем тело запроса

	// Ошибка в декодере
	err := decoder.Decode(&book) // Передаем адресс на книгу
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = err.Error()

		respJson, _ := json.Marshal(response)
		_, _ = w.Write(respJson)

		return
	}

	// Ошибка при обновлении книги
	err = bookStore.UpdateBook(book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = err.Error()

		respJson, _ := json.Marshal(response)
		_, _ = w.Write(respJson)

		return
	}

	w.WriteHeader(http.StatusOK)

	response = Response{
		Message: fmt.Sprintf("Book with id:%s updated is SUCCESS!", book.Id),
	}
	respJson, _ := json.Marshal(response)

	_, _ = w.Write(respJson)
}

func handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var response Response

	// Считываем книгу
	decoder := json.NewDecoder(r.Body) // Декодируем тело запроса

	// Ошибка в декодере
	err := decoder.Decode(&book) // Передаем адресс на книгу
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = err.Error()

		respJson, _ := json.Marshal(response)
		_, _ = w.Write(respJson)

		return
	}

	// Ошибка при удалении книги
	err = bookStore.DeleteBook(book.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = err.Error()

		respJson, _ := json.Marshal(response)
		_, _ = w.Write(respJson)

		return
	}

	w.WriteHeader(http.StatusOK)

	response = Response{
		Message: fmt.Sprintf("Book with id:%s deleted is SUCCESS!", book.Id),
	}
	respJson, _ := json.Marshal(response)

	_, _ = w.Write(respJson)
}

///////////////////////////////////////////////////////////////////

func (s BookStore) FindBookById(id string) *Book {
	for _, book := range s.books {
		if book.Id == id {
			return &book
		}
	}

	return nil
}

func (s BookStore) GetBooks() []Book {
	return s.books
}

func (s BookStore) GetBook(id string) *Book {
	for _, book := range s.books {
		if book.Id == id {
			return &book
		}
	}

	return nil
}

func (s *BookStore) AddBook(book Book) error {
	for _, bk := range s.books {
		if bk.Id == book.Id {
			return errors.New(fmt.Sprintf("Book with Id:%s is isset!", book.Id))
		}
	}

	s.books = append(s.books, book)

	return nil
}

func (s *BookStore) UpdateBook(book Book) error {
	for i, bk := range s.books {
		if bk.Id == book.Id {
			s.books[i] = book
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Book with Id:%s not found!", book.Id))
}

func (s *BookStore) DeleteBook(id string) error {
	for i, bk := range s.books {
		if bk.Id == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Book with Id:%s not found!", id))
}
