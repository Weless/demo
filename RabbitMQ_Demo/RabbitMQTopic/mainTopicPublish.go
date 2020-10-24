package main

import (
	"demo/RabbitMQ_Demo/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	imOne := RabbitMQ.NewRabbitMQTopic("exIM2", "im.topic.one")
	imTwo := RabbitMQ.NewRabbitMQTopic("exIM2", "im.topic.two")
	for i := 0; i <= 10; i++ {
		imOne.PublishTopic("hello one" + strconv.Itoa(i))
		imTwo.PublishTopic("hello two" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
