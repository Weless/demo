package main

import "demo/RabbitMQ_Demo/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("newProduct")
	rabbitmq.ReceiveSub()
}
