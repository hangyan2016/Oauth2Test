/**
 * @Author: Guo PingFei
 * @Description:
 * @File:  receive
 * @Version: 1.0.0
 * @Date: 2019/12/12 10:50
 */


package main

import (
"log"

"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Agent", // name
		true,    // durable
		true,    // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	err = ch.QueueBind("Agent", "Agent", "Test_01", false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Println(d.Type)
			log.Println(d.MessageId)
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}

