package main

import "demo/RabbitMQ_Demo/RabbitMQ"

func main() {
	imTwo := RabbitMQ.NewRabbitMQTopic("exIM2", "im.*.two")
	imTwo.ReceiveTopic()
}
