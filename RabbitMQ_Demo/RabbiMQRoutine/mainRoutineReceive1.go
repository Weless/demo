package main

import "demo/RabbitMQ_Demo/RabbitMQ"

func main() {
	imOne := RabbitMQ.NewRabbitMQRoutine("exIM", "im_one")
	imOne.ReceiveRouting()
}
