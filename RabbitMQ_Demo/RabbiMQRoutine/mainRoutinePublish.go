package main

import (
	"demo/RabbitMQ_Demo/RabbitMQ"
	"fmt"
	"time"
)

func main() {
	imOne := RabbitMQ.NewRabbitMQRoutine("exIM", "im_one")
	imTwo := RabbitMQ.NewRabbitMQRoutine("exIM", "im_two")

	for i := 0; i <= 10; i++ {
		imOne.PublishRouting("hello one!")
		imTwo.PublishRouting("hello two!")
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
