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

	// Consume messages from the Redis queue
	consumeMessages(client)
}

func consumeMessages(client *redis.Client) {
	for {
		result, err := client.BLPop(client.Context(), 0, "myqueue").Result()
		if err != nil {
			log.Printf("Failed to consume message: %v", err)
		} else {
			message := result[1]
			fmt.Println("Worker consumed message:", message)

			// Perform 50% tasks with the consumed message
			status := performTaskOne(message)

			// Report task status as 50% done back to the server
			reportStatus(client, message, status)

			// Perform remaining 50% tasks with the consumed message
			status = performTaskTwo(message)

			// Report task status as 100% back to the server
			reportStatus(client, message, status)

			fmt.Println("Worker finished performing Task:", message)
		}

		time.Sleep(time.Second) // Sleep for a second between consuming messages
	}
}

func performTaskOne(message string) string {
	// Simulating task processing time
	time.Sleep(3 * time.Second)

	// Return a dummy task status
	return "Completed 50% Task"
}

func performTaskTwo(message string) string {
	// Simulating task processing time
	time.Sleep(3 * time.Second)

	// Return a dummy task status
	return "Completed 100% Task"
}

func reportStatus(client *redis.Client, message, status string) {
	err := client.RPush(client.Context(), "statusqueue", fmt.Sprintf("%s: %s", message, status)).Err()
	if err != nil {
		log.Printf("Failed to report task status: %v", err)
	} else {
		fmt.Printf("Reported task status: %s - %s\n", message, status)
	}
}
