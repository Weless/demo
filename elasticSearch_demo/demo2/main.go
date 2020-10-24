package main

import (
	"context"
	"elasticsearch_demo/demo2/AppInit"
	"elasticsearch_demo/demo2/Models"
	"fmt"
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

		for _, book := range bookList {
			ctx := context.Background()
			rsp, err := AppInit.NewESClient().Index().Index("books").Id(strconv.Itoa(book.BookID)).
				BodyJson(book).Do(ctx)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rsp)
			}
		}

		break
	}
}
