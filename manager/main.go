package main

import (
	"context"
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

	// Receive worker queue info from workers
	fmt.Println(" ----- HANDSHAKE WITH WORKERS ----- ")
	workerQueues := handshakewithAllWorkers(client)
	fmt.Println("Received all Worker queues as : ", workerQueues)

	fmt.Println(" ----- CREATE AND PUBLISH JOBS FOR ALL WORKERS TO CONSUME ----- ")
	// Create and Publish multiple jobs to the queue
	createAndPublishJobs(client)

	fmt.Println(" ----- SUBSCRIBE TO JOB_STATUS_TOPIC TO GET JOB EXECUTION STATUS FOR EVERY WORKER ----- ")
	// Subscribe to the job status topic
	subscribeJobStatusTopic(client)

	fmt.Println(" ----- SEND STOP JOB EXECUTION MESSAGE TO ALL WORKERS QUEUE ----- ")
	//send stop message to all workers queue
	stopJobExecutionOnAllWorkers(client, workerQueues)

	// Block indefinitely
	select {}
}
func handshakewithAllWorkers(client *redis.Client) []string {
	workerQueueSubscriber := client.Subscribe(client.Context(), "worker_queue_info")
	_, err := workerQueueSubscriber.Receive(client.Context())
	if err != nil {
		log.Fatalf("Failed to subscribe to worker queue info: %v", err)
	}
	// Wait for all workers to publish their queue info
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	infoCh := workerQueueSubscriber.Channel()
	workerQueueList := make([]string, 0)
	for {
		select {
		case msg := <-infoCh:
			workerQueueList = append(workerQueueList, msg.Payload)
		case <-ctx.Done():
			fmt.Println("Stopping subscription after getting all workers queue info")
			err := workerQueueSubscriber.Unsubscribe(context.Background(), "worker_queue_info")
			if err != nil {
				log.Printf("Failed to unsubscribe from worker_queue_info: %v", err)
			}
			return workerQueueList
		}
	}
}
func receiveWorkerQueueInfo(channel <-chan *redis.Message, workerQueueCh chan<- string) {
	for msg := range channel {
		workerQueue := msg.Payload
		fmt.Println("Received worker queue info:", workerQueue)
	}
}
func createAndPublishJobs(client *redis.Client) {
	for i := 1; i <= 5; i++ {
		job := "Job " + strconv.Itoa(i)
		err := client.LPush(client.Context(), "jobs_queue", job).Err()
		if err != nil {
			log.Printf("Failed to publish: %v", err)
		} else {
			fmt.Println("Published:", job)
		}
		time.Sleep(time.Second)
	}
}
func subscribeJobStatusTopic(client *redis.Client) {
	jobStatusTopicSubscriber := client.Subscribe(client.Context(), "job_status_topic")
	_, err := jobStatusTopicSubscriber.Receive(client.Context())
	if err != nil {
		log.Fatalf("Failed to subscribe to job status topic: %v", err)
	}
	go receiveJobStatusUpdates(jobStatusTopicSubscriber.Channel())
}
func stopJobExecutionOnAllWorkers(client *redis.Client, workerQueues []string) {
	for _, workerQueue := range workerQueues {
		fmt.Println("Sending stop message to worker queue:", workerQueue)
		err := client.LPush(client.Context(), workerQueue, "stop").Err()
		if err != nil {
			log.Printf("Failed to publish: %v", err)
		} else {
			fmt.Println("Published stop message to worker queue:", workerQueue)
		}
	}
}

func receiveJobStatusUpdates(channel <-chan *redis.Message) {
	for msg := range channel {
		jobStatus := msg.Payload
		fmt.Println("Received job status update:", jobStatus)
	}
}
