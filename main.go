package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "2", Title: "The River and The Source", Author: "Millicent Ogola", Quantity: 3},
	{ID: "4", Title: "Life In Crime", Author: "John Kiriamiti", Quantity: 5},
	{ID: "1", Title: "The Bourne Identity", Author: "Unknown", Quantity: 3},
	{ID: "3", Title: "In search of lost time", Author: "Marsel Proust", Quantity: 3},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func bookById(c *gin.Context) {
	id := c.Param("id")

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func confirmBook(c *gin.Context) *book {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "missing query id parameter."})
		return nil
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found."})
		return nil
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not available."})
		return nil
	}

	return book
}

func bookCheckout(c *gin.Context) {
	book := confirmBook(c)

	book.Quantity -= 1

	c.IndentedJSON(http.StatusOK, book)
}

func bookReturn(c *gin.Context) {
	book := confirmBook(c)

	book.Quantity += 1

	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", bookCheckout)
	router.PATCH("/return", bookReturn)
	router.Run("localhost:5000")
}
