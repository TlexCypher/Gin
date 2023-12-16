package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json: "title"`
	Author   string `json: author`
	Quantity int    `json : "quantity"`
}

var books = []book{
	{ID: "1", Title: "Nvim is awesome", Author: "funasoul", Quantity: 123},
	{ID: "2", Title: "Emacs is awesome", Author: "tourist", Quantity: 321},
	{ID: "3", Title: "Vscode is awesome", Author: "tokkuman", Quantity: 111},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var b book
	if err := c.BindJSON(&b); err != nil {
		return
	}
	books = append(books, b)
	c.IndentedJSON(http.StatusCreated, b)
}

func fetchBookById(c *gin.Context) {
	id := c.Param("id")
	b, err := findBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"msg": "Error has been occured. Such book was not found.",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, b)
}

func findBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Such book was not found.")
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Missing id query parameter."})
		return
	}
	b, err := findBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"msg": "Error has been occured. Suck kind of book was not found.",
		})
		return
	}
	if b.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"msg": "The books quantity is less than zero.",
		})
		return
	}
	b.Quantity -= 1

	c.IndentedJSON(http.StatusOK, gin.H{
		"msg": "The book was checkouted correctly. The current quantity of the book is " + strconv.Itoa(b.Quantity),
	})
}

func main() {
	r := gin.Default()
	r.GET("/books", getBooks)
	r.GET("/books/:id", fetchBookById)
	r.PATCH("/checkout", checkoutBook)
	r.POST("/books", createBook)
	r.Run("localhost:8080")
}
