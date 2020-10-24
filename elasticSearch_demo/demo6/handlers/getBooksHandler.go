package handlers

import (
	"elasticsearch_demo/demo5/AppInit"
	"elasticsearch_demo/demo6/Models"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"reflect"
)

func MapToBooks(rsp *elastic.SearchResult) []*Models.Books {
	books := []*Models.Books{}
	var t *Models.Books
	for _, item := range rsp.Each(reflect.TypeOf(t)) {
		books = append(books, item.(*Models.Books))
	}
	return books
}

func GetBooksHandler(c *gin.Context) {
	rsp, err := AppInit.NewESClient().Search().Index("books").Do(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	} else {
		c.JSON(200, gin.H{"result": MapToBooks(rsp)})
	}
}
