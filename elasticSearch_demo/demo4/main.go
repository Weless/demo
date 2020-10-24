package main

import (
	"context"
	"elasticsearch_demo/demo2/AppInit"
	"elasticsearch_demo/demo2/Models"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strconv"
	"sync"
)

func main() {
	page := 1
	pageSize := 500
	wg := sync.WaitGroup{}
	for {
		bookList := Models.BookList{}
		db := AppInit.GetDB().Select("book_id,book_name,book_intr,book_price1,book_price2,book_author,book_press,book_kind " +
			",if(book_date='','1970-01-01',ltrim(book_date)) as book_date").
			Order("book_id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&bookList)
		if db.Error != nil || len(bookList) == 0 {
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
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
		}()
		page = page + 1
	}
	wg.Wait()
}
