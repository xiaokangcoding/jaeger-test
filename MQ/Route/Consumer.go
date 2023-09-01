package main

import (
	"github.com/Henry-jk/jaeger-test/MQ/Route/util"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"Direct-Exchange", // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	queues := []string{"queue1", "queue2", "queue3","queue3"}
	routingKeys := []string{"email", "sms", "weixin","email"}

	forever := make(chan bool)

	for i, queueName := range queues {
		q, err := ch.QueueDeclare(
			queueName, // name
			true ,    // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		util.FailOnError(err, "Failed to declare a queue")

		err = ch.QueueBind(
			q.Name,         // queue name
			routingKeys[i], // routing key
			"Direct-Exchange",   // exchange
			false,
			nil,
		)
		util.FailOnError(err, "Failed to bind a queue")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		util.FailOnError(err, "Failed to register a consumer")

		go func(queueName string, msgs <-chan amqp.Delivery) {
			for d := range msgs {
				log.Printf("Received a message from queue %s: %s", queueName, d.Body)
			}
		}(queueName, msgs)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//func main() {
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//	util.FailOnError(err, "Failed to connect to RabbitMQ")
//	defer conn.Close()
//
//	ch, err := conn.Channel()
//	util.FailOnError(err, "Failed to open a channel")
//	defer ch.Close()
//
//	err = ch.ExchangeDeclare(
//		"carl-route", // name
//		"direct",   // type
//		true,       // durable
//		false,      // auto-deleted
//		false,      // internal
//		false,      // no-wait
//		nil,        // arguments
//	)
//	util.FailOnError(err, "Failed to declare an exchange")
//
//	queues := []string{"hello1", "hello2", "simple_queue"}
//	routingKeys := []string{"email", "email", "weixin"}
//
//	for i, queueName := range queues {
//		q, err := ch.QueueDeclare(
//			queueName, // name
//			false,     // durable
//			false,     // delete when unused
//			false,     // exclusive
//			false,     // no-wait
//			nil,       // arguments
//		)
//		util.FailOnError(err, "Failed to declare a queue")
//
//		err = ch.QueueBind(
//			q.Name,         // queue name
//			routingKeys[i], // routing key
//			"carl-route",     // exchange
//			false,
//			nil,
//		)
//		util.FailOnError(err, "Failed to bind a queue")
//
//		msgs, err := ch.Consume(
//			q.Name, // queue
//			"",     // consumer
//			true,   // auto-ack
//			false,  // exclusive
//			false,  // no-local
//			false,  // no-wait
//			nil,    // args
//		)
//		util.FailOnError(err, "Failed to register a consumer")
//
//		forever := make(chan bool)
//
//		go func() {
//			for d := range msgs {
//				log.Printf("Received a message from queue %s: %s", queueName, d.Body)
//			}
//		}()
//
//		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
//		<-forever
//	}
//}

//var queueName=[]string{"simple-queue","hello1","hello2"}
//
//func main() {
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//	util.FailOnError(err, "Failed to connect to RabbitMQ")
//	defer conn.Close()
//
//	ch, err := conn.Channel()
//	util.FailOnError(err, "Failed to open a channel")
//	defer ch.Close()
//
//	for _, name := range queueName {
//
//	}
//	err = ch.ExchangeDeclare(
//		"carl-route", // name
//		"direct",      // type
//		true,          // durable
//		false,         // auto-deleted
//		false,         // internal
//		false,         // no-wait
//		nil,           // arguments
//	)
//	util.FailOnError(err, "Failed to declare an exchange")
//
//	q, err := ch.QueueDeclare(
//		"simple-queue",    // name
//		false, // durable
//		false, // delete when unused
//		true,  // exclusive
//		false, // no-wait
//		nil,   // arguments
//	)
//	util.FailOnError(err, "Failed to declare a queue")
//
//	if len(os.Args) < 2 {
//		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
//		os.Exit(0)
//	}
//	for _, s := range os.Args[1:] {
//		log.Printf("Binding queue %s to exchange %s with routing key %s",
//			q.Name, "logs_direct", s)
//		err = ch.QueueBind(
//			q.Name,        // queue name
//			s,             // routing key
//			"fanout-exchange", // exchange
//			false,
//			nil)
//		util.FailOnError(err, "Failed to bind a queue")
//	}
//
//	msgs, err := ch.Consume(
//		q.Name, // queue
//		"",     // consumer
//		true,   // auto ack
//		false,  // exclusive
//		false,  // no local
//		false,  // no wait
//		nil,    // args
//	)
//	util.FailOnError(err, "Failed to register a consumer")
//
//	forever := make(chan bool)
//
//	go func() {
//		for d := range msgs {
//			log.Printf(" [x] %s", d.Body)
//		}
//	}()
//
//	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
//	<-forever
//}
