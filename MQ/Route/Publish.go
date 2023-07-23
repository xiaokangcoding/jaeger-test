package main

import (
	"github.com/Henry-jk/jaeger-test/MQ/Route/util"
	"log"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"carl-route",   // name
		"direct",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	routingKeys := []string{"email", "weixin"}

	for _, routingKey := range routingKeys {
		body := "Hello World!"
		err = ch.Publish(
			"carl-route",     // exchange
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
//		"direct",      // type
//		true,          // durable
//		false,         // auto-deleted
//		false,         // internal
//		false,         // no-wait
//		nil,           // arguments
//	)
//	util.FailOnError(err, "Failed to declare an exchange")
//
//	body := bodyFrom(os.Args)
//	err = ch.Publish(
//		"carl-route",         // exchange
//		severityFrom(os.Args), // routing key
//		false, // mandatory
//		false, // immediate
//		amqp.Publishing{
//			ContentType: "text/plain",
//			Body:        []byte(body),
//		})
//	util.FailOnError(err, "Failed to publish a message")
//
//	log.Printf(" [x] Sent %s", body)
//}
//
//func bodyFrom(args []string) string {
//	var s string
//	if (len(args) < 3) || os.Args[2] == "" {
//		s = "hello"
//	} else {
//		s = strings.Join(args[2:], " ")
//	}
//	return s
//}
//
//func severityFrom(args []string) string {
//	var s string
//	if (len(args) < 2) || os.Args[1] == "" {
//		s = "info"
//	} else {
//		s = os.Args[1]
//	}
//	return s
//}