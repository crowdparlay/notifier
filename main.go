package main

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Notification struct {
	UserId int `json:"userId"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error while connecting to AMQP: %s", err)
	}
	defer func(conn *amqp.Connection) {
		_ = conn.Close()
	}(conn)

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error while connecting to channel: %s", err)
	}
	defer func(channel *amqp.Channel) {
		_ = channel.Close()
	}(channel)

	notifications, err := channel.Consume(
		"notifications",
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	log.Println("Connected to channel")

	loop := make(chan bool)

	go func() {
		for notification := range notifications {
			log.Printf("Get message: %s \n", notification.Body)

			var m Notification

			err = json.Unmarshal(notification.Body, &m)

			if err != nil {
				log.Fatalf("Fucked up: %s", err)
			}

			log.Printf("JSON: %+v", m)
		}
	}()

	<-loop
}
