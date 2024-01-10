package internal

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type Notification struct {
	UserId int `json:"userId"`
}

func ServeWS(c *gin.Context, upgrader *websocket.Upgrader, notifications <-chan amqp091.Delivery) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("Error while upgrading WS connection: %s", err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	for notification := range notifications {
		log.Printf("Get message: %s \n", notification.Body)

		var m Notification

		err = json.Unmarshal(notification.Body, &m)

		if err != nil {
			log.Fatalf("Fucked up: %s", err)
		}

		err := conn.WriteJSON(m)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}
