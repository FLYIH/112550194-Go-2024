package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Book struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Pages int    `json:"pages"`
}

var bookshelf = []Book{
	{ID: 1, Name: "Blue Bird", Pages: 500},
}

var nextID = 2

func sendError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"message": message})
}

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, bookshelf)
}

func getBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sendError(c, http.StatusBadRequest, "Invalid book ID")
		return
	}

	for _, book := range bookshelf {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}
	sendError(c, http.StatusNotFound, "book not found")
}

func addBook(c *gin.Context) {
	var newBook Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if newBook.Pages < 1 {
		sendError(c, http.StatusBadRequest, "Pages must be greater than 0")
		return
	}

	for _, book := range bookshelf {
		if book.Name == newBook.Name {
			sendError(c, http.StatusConflict, "duplicate book name")
			return
		}
	}

	newBook.ID = nextID
	nextID++
	bookshelf = append(bookshelf, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func deleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sendError(c, http.StatusBadRequest, "Invalid book ID")
		return
	}

	for i, book := range bookshelf {
		if book.ID == id {
			bookshelf = append(bookshelf[:i], bookshelf[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.Status(http.StatusNoContent)
}

func updateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sendError(c, http.StatusBadRequest, "Invalid book ID")
		return
	}

	for i, book := range bookshelf {
		if book.ID == id {
			var updatedBook Book
			if err := c.ShouldBindJSON(&updatedBook); err != nil {
				sendError(c, http.StatusBadRequest, err.Error())
				return
			}

			if updatedBook.Pages < 1 {
				sendError(c, http.StatusBadRequest, "Pages must be greater than 0")
				return
			}

			for _, b := range bookshelf {
				if b.Name == updatedBook.Name && b.ID != id {
					sendError(c, http.StatusConflict, "duplicate book name")
					return
				}
			}

			updatedBook.ID = id
			bookshelf[i] = updatedBook
			c.JSON(http.StatusOK, updatedBook)
			return
		}
	}
	sendError(c, http.StatusNotFound, "book not found")
}

func main() {
	r := gin.Default()
	r.RedirectFixedPath = true

	r.GET("/bookshelf", getBooks)
	r.GET("/bookshelf/:id", getBook)
	r.POST("/bookshelf", addBook)
	r.PUT("/bookshelf/:id", updateBook)
	r.DELETE("/bookshelf/:id", deleteBook)

	err := r.Run(":8087")
	if err != nil {
		return
	}
}
