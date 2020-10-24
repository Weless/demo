package handlers

import (
	"elasticsearch_demo/demo7/AppInit"
	"elasticsearch_demo/demo7/Models"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strings"
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

func GetBooksByPress(c *gin.Context) {
	press, _ := c.Params.Get("press")
	termQuery := elastic.NewTermQuery("BookPress", press)
	rsp, err := AppInit.NewESClient().Search().Query(termQuery).Index("books").Do(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	} else {
		c.JSON(200, gin.H{"result": MapToBooks(rsp)})
	}
}

func GetBooksByPresses(c *gin.Context) {
	press, _ := c.Params.Get("press")
	list := strings.Split(press, ",")
	pressList := []interface{}{}
	for _, p := range list {
		pressList = append(pressList, p)
	}
	termQuery := elastic.NewTermsQuery("BookPress", pressList...)
	rsp, err := AppInit.NewESClient().Search().Query(termQuery).Index("books").Do(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	} else {
		c.JSON(200, gin.H{"result": MapToBooks(rsp)})
	}
}
