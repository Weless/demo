package main

import (
	"demo/RabbitMQ_Demo/RabbitMQ"
	"fmt"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imSimple")
	rabbitmq.PublishSimple("hello joey")
	fmt.Println("send success")
}
