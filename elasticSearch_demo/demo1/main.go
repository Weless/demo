package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

func main() {
	client, err := elastic.NewClient(
		// es地址，如果是集群，可加多个
		elastic.SetURL("http://127.0.0.1:9200/"),
		// 需要加上这个配置，否则找不到（请求连接返回的是docker的ip地址）？？ 若es和本机能够连通，设置为true
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	mapping, err := client.GetMapping().Index("news").Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(mapping)
}
