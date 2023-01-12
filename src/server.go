package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloResponse struct {
	Title            string
	ShortDescription string
}

func main() {
	router := gin.Default()
	router.GET("/hello", helloWorld)

	router.Run("localhost:8080")
}

func helloWorld(c *gin.Context) {
	m := HelloResponse{"Test", "Hello"}
	c.IndentedJSON(http.StatusOK, m)
}
