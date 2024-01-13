package main

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
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

	_, err = channel.QueueDeclare(
		"notifications",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Error while declareting queue: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := Notification{
		UserId: 2,
	}

	bytes, _ := json.Marshal(body)

	if err := channel.PublishWithContext(ctx, "", "notifications", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bytes,
	}); err != nil {
		log.Fatalf("Error while publishing message: %s", err)
	}
}
