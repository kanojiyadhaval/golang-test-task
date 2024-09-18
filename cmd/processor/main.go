package main

import (
	"context"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

var ctx = context.Background()

func main() {
	conn, err := amqp.Dial("amqp://user:password@localhost:7001/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
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
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis container address
		Password: "",           // No password set by default
		DB:       0,            // Use default DB
	})
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Printf("Redis connected: %v", pong)
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			processMessage(redisClient, string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func processMessage(redisClient *redis.Client, body string) {
	parts := strings.Split(body, ",")
	if len(parts) != 3 {
		log.Printf("Invalid message format: %s", body)
		return
	}
	sender := parts[0]
	receiver := parts[1]
	message := parts[2]

	key := sender + "_" + receiver
	log.Printf("Processing message: sender=%s, receiver=%s, message=%s", sender, receiver, message)

	redisClient.LPush(ctx, key, message)
}
