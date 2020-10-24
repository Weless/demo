package main

import (
	"context"
	"elasticsearch_demo/demo2/AppInit"
	"elasticsearch_demo/demo2/Models"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strconv"
)

func main() {
	page := 1
	pageSize := 50
	for {
		bookList := Models.BookList{}
		db := AppInit.GetDB().Order("book_id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&bookList)
		if db.Error != nil || len(bookList) == 0 {
			break
		}

		client := AppInit.NewESClient()
		bulk := client.Bulk()
		for _, book := range bookList {
			req := elastic.NewBulkIndexRequest().Index("books").Id(strconv.Itoa(book.BookID)).Doc(book)
			bulk.Add(req)
		}
		rsp, err := bulk.Do(context.Background())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(rsp)
		}
		break
	}
}
