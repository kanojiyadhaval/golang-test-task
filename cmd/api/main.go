package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Message struct {
	Sender   string `json:"sender" binding:"required"`
	Receiver string `json:"receiver" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

func main() {
	r := gin.Default()

	r.POST("/message", func(c *gin.Context) {
		var msg Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := pushToQueue(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error sending message to RabbitMQ"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Message sent successfully"})
	})

	r.Run(":8080")
}

func pushToQueue(msg Message) error {
	conn, err := amqp.Dial("amqp://user:password@localhost:7001/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"message_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return err
	}

	body := msg.Sender + "," + msg.Receiver + "," + msg.Message
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
		return err
	}

	return nil
}
