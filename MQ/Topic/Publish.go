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
		"topic-exchange",   // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	//routingKeys := []string{"com.#", "*.course.*","#.order.#"}
	routingKeys := []string{"com.","com.xxx" ,"com.XXX.VVV","aaa.course.ddd",".order.ccc"}
	for _, routingKey := range routingKeys {
		body := "Hello World!"
		err = ch.Publish(
			"topic-exchange",     // exchange
			routingKey, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		util.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)
	}
}
