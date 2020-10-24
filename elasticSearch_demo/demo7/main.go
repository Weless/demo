package main

import (
	"elasticsearch_demo/demo7/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	g := r.Group("/books")
	{
		g.GET("", handlers.GetBooksHandler)
		g.GET("/press/:press", handlers.GetBooksByPress)
		g.GET("/presses/:press", handlers.GetBooksByPresses)
	}
	r.Run(":8080")
}
