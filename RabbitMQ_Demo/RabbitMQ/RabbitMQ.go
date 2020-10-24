package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// url 格式 amqp//账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://joey:qwe123@127.0.0.1:5672/admin"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// Key
	Key string
	// 连接信息
	MqUrl string
}

func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
	}
	rabbitmq.MqUrl = MQURL
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "创建连接错误!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

// 断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// 简单模式Step：1.创建简单模式下rabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

// 简单模式Step：2.简单模式下的生产代码
func (r *RabbitMQ) PublishSimple(message string) {
	// 1. 申请队列
	// 队列如果不存在，队列将会被创建；如果队列存在，将会跳过，
	// 优点：确保生产者的消息能够发送到队列里
	_, err := r.channel.QueueDeclare(
		// 队列名称
		r.QueueName,
		// durable：控制消息是否持久化，服务器重启消息就没了
		false,
		// autoDelete：是否自动删除，当最后一个消费者断开连接后，是否把消息从队列中删除
		false,
		// exclusive： 是否具有排他性，当为true的时候就创建了一个仅为自己可见的队列，其他用户无法访问
		false,
		// noWait:是否阻塞，发送消息后是否等待服务器的响应
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		fmt.Printf("Declare queue failed, err:%s\n", err)
	}
	// 2. 发送消息到队列中
	r.channel.Publish(
		// exchange 为空的时候，为默认的交换机
		r.Exchange,
		// queueName 队列名
		r.QueueName,
		// mandatory: 如果是true，会根据自身的exchange类型和routekey规则，判断是否能找到符合条件的队列，如果找不到队列，会把消息返还给消费者
		false,
		// immediate：如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息发还给发送者
		false,
		// msg:发送的消息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 简单方式Step：3.简单模式下消费代码
func (r *RabbitMQ) ConsumeSimple() {
	// 1. 申请队列
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Printf("Declare queue failed, err:%s\n", err)
	}

	// 接收消息
	msgs, err := r.channel.Consume(
		// queue: 队列
		r.QueueName,
		// consumer: 用来区分多个消费者
		"",
		// autoAck: 是否自动应答，是否主动告诉rabbitMQ消息已经消费完了；如果为false，需要手动实现回调函数
		true,
		// exclusive： 是否具有排他性；
		false,
		// noLocal: 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 队列消费是否阻塞
		false,
		nil)
	if err != nil {
		fmt.Printf("Consume failed, err:%s\n", err)
	}

	forever := make(chan bool)
	// 处理消息
	go func() {
		for d := range msgs {
			// 实现逻辑函数
			log.Printf("Received a messagez: %s\n", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages, to exit press CTRL+C\n")
	<-forever
}

// 订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	// 订阅模式key必须为空
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "failed to connect to rabbitmq")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型，订阅模式下需要设置为fanout,即广播类型
		"fanout",
		true,
		false,
		// true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的相互绑定
		false,
		false,
		nil)
	r.failOnErr(err, "fail to declare an exchange")

	// 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 订阅模式消费端代码
func (r *RabbitMQ) ReceiveSub() {
	// 试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		// YES表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an exchange")

	// 试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		// 生成一个随机的队列名
		"",
		false,
		false,
		true,
		false,
		nil)
	r.failOnErr(err, "failed to declare a queue")

	// 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		// 在pub/sub模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil)
	r.failOnErr(err, "failed to bind queue to exchange")

	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to consume the messages")

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a messge: %s", d.Body)
		}
	}()
	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}

// Routine模式创建RabbitMQ实例
func NewRabbitMQRoutine(exchangeName string, routineKey string) *RabbitMQ {
	// key 参数需要传入routinekey
	rabbitmq := NewRabbitMQ("", exchangeName, routineKey)
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "failed to connect to rabbitmq")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// Routine模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 要改成direct
		"direct",
		true,
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an exchange")

	// 2. 发送消息
	err = r.channel.Publish(
		r.Exchange,
		// 要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

// Routine模式接收消息
func (r *RabbitMQ) ReceiveRouting() {
	// 1. 试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型
		"direct",
		true,
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an exchange")

	// 2.试探性创建队列，注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	r.failOnErr(err, "failed to declare a queue")

	// 3. 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil)
	r.failOnErr(err, "failed to bind the queue to the exchange")
	// 4. 消费消息

	messages, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)
	r.failOnErr(err, "failed to consume")

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a messgae: %s", d.Body)
		}
	}()
	fmt.Println("退出请按 CTRL+C")
	<-forever
}

// 话题模式创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName, routineKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routineKey)
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "failed to connect to rabbitmq")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string) {
	// "*"用来匹配一个单词，"#"用来匹配多个单词（可以是零个）
	// imooc.* 匹配 imooc.hello, imooc.# 可以匹配imooc.hello.one

	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an exchange")

	err = r.channel.Publish(r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 话题模式接收消息
func (r *RabbitMQ) ReceiveTopic() {
	err := r.channel.ExchangeDeclare(r.Exchange, "topic", true, false, false, false, nil)
	r.failOnErr(err, "failed to declare an exchange")

	queue, err := r.channel.QueueDeclare("", false, false, true, false, nil)
	r.failOnErr(err, "failed to declare a queue")

	err = r.channel.QueueBind(queue.Name, r.Key, r.Exchange, false, nil)
	r.failOnErr(err, "failed to bind the queue to the exchange")

	messages, err := r.channel.Consume(queue.Name, "", true, false, false, false, nil)
	r.failOnErr(err, "failed to consume")

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a messgae: %s", d.Body)
		}
	}()
	fmt.Println("退出请按 CTRL+C")
	<-forever

}
