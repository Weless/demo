package handlers

import (
	"elasticsearch_demo/demo5/AppInit"
	"github.com/gin-gonic/gin"
)

func GetBooksHandler(c *gin.Context) {
	rsp, err := AppInit.NewESClient().Search().Index("books").Do(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	} else {
		c.JSON(200, gin.H{"result": rsp.Hits.Hits})
	}
}
