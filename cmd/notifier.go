package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"notifier/internal"
	"os"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	url := os.Getenv("AMQP_CONNECTION_URL")
	conn, err := amqp.Dial(url)
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

	router := gin.Default()

	router.GET("/notifications/:id", func(c *gin.Context) {
		internal.ServeWS(c, upgrader, notifications)
	})

	err = router.Run(":8000")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
