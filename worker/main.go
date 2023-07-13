package main

import (
	"fmt"
	"log"
	"strconv"
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

	// Wait for the manager to subscribe to the worker queue info before publishing
	time.Sleep(5 * time.Second)

	// Publish the worker Queue information to the server
	workerQueue := publisWorkerQueueInfo(client)

	go consumeWorkerQueue(client, workerQueue)

	// Consume jobs from the jobs queue
	consumeJobs(client, "jobs_queue")
}

func publisWorkerQueueInfo(client *redis.Client) string {
	workerQueue := "worker_queue_" + generateUniqueID()

	err := client.Publish(client.Context(), "worker_queue_info", workerQueue).Err()
	if err != nil {
		log.Fatalf("Failed to publish worker queue info: %v", err)
	}
	return workerQueue
}

func consumeWorkerQueue(client *redis.Client, workerQueue string) {
	for {
		// Consume worker queue for any stop message from the server
		message, err := client.RPop(client.Context(), workerQueue).Result()
		if err != nil {
			//log.Printf("Failed to consume worker queue: %v", err)
			continue
		}
		fmt.Printf("Received message [%s] on worker queue [%s] to stop job execution:\n", message, workerQueue)
	}
}

func consumeJobs(client *redis.Client, jobsQueue string) {
	for {
		// Consume job from the jobs queue
		job, err := client.RPop(client.Context(), jobsQueue).Result()
		if err != nil {
			//log.Printf("Failed to consume job: %v", err)
			continue
		}

		if job == "" {
			fmt.Println("No more jobs in the queue. Exiting...")
			return
		}

		fmt.Println("Consumed job:", job)

		// Perform job execution
		time.Sleep(5 * time.Second) // Simulating job processing time

		// Publish job status update
		jobStatus := fmt.Sprintf("%s executed successfully", job)
		err = client.Publish(client.Context(), "job_status_topic", jobStatus).Err()
		if err != nil {
			log.Printf("Failed to publish job status update: %v", err)
		}
	}
}

func generateUniqueID() string {
	// Generate a unique ID based on the current timestamp
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
