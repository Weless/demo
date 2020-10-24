package main

import (
	"demo/RabbitMQ_Demo/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imSimple")
	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("hello world" + strconv.Itoa(i))
		time.Sleep(2 * time.Second)
		fmt.Println(i)
	}
}
