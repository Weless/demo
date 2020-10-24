package main

import "demo/RabbitMQ_Demo/RabbitMQ"

func main() {
	imOne := RabbitMQ.NewRabbitMQTopic("exIM2", "#")
	imOne.ReceiveTopic()
}
