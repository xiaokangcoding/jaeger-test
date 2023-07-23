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
		"topic-exchange", // name
		"topic",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	queues := []string{"hello1", "hello2", "simple-queue"}
	routingKeys := []string{"com.#", "*.course.*","#.order.#"}

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
			"topic-exchange",   // exchange
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