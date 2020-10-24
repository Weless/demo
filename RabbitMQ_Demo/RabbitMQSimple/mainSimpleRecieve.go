package main

import "demo/RabbitMQ_Demo/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imSimple")
	rabbitmq.ConsumeSimple()
}
