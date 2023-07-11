package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Use the service name defined in the docker-compose.yml file
		Password: "",           // Replace with your Redis server password, if any
		DB:       0,            // Replace with the desired Redis database index
	})

	// Ping Redis to check the connection
	pong, err := client.Ping(client.Context()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)

	// Publish messages to a Redis queue
	publishMessages(client)

	//Consume messages from the Redis status queue
	consumeMessages(client)
}

func publishMessages(client *redis.Client) {
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Message %d", i)

		err := client.RPush(client.Context(), "myqueue", message).Err()
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		} else {
			fmt.Println("Published message:", message)
		}

		time.Sleep(time.Second) // Sleep for a second between publishing messages
	}
}

func consumeMessages(client *redis.Client) {
	for {
		result, err := client.BLPop(client.Context(), 0, "statusqueue").Result()
		if err != nil {
			log.Printf("Failed to consume message: %v", err)
		} else {
			message := result[1]
			fmt.Println("Manager consumed message regarding worker status on Task:", message)
		}

		time.Sleep(time.Second) // Sleep for a second between consuming messages
	}
}
