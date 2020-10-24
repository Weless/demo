package main

import (
	"elasticsearch_demo/demo5/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	g := r.Group("/books")
	{
		g.GET("", handlers.GetBooksHandler)
	}
	r.Run(":8080")
}
