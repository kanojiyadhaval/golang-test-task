package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	r := gin.Default()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	r.GET("/message/list", func(c *gin.Context) {
		sender := c.Query("sender")
		receiver := c.Query("receiver")

		if sender == "" || receiver == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sender and receiver are required"})
			return
		}

		key := sender + "_" + receiver
		messages, err := redisClient.LRange(ctx, key, 0, -1).Result()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"messages": messages})
	})

	r.Run(":8081")
}
