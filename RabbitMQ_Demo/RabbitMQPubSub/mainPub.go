package main

import (
	"demo/RabbitMQ_Demo/RabbitMQ"
	"log"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("newProduct")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("订阅模式生产第" + strconv.Itoa(i) + "条数据")
		log.Printf("订阅模式生产第%d条数据", i)
		time.Sleep(1 * time.Second)
	}
}
