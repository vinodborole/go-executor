# Go-Executor

This project demonstrates a simple implementation of a Redis queue system using Go programs and Docker containers. It consists of a manager program that publishes messages to a Redis queue and multiple worker programs that consume the messages and perform tasks concurrently. Additionally, the workers report their task status back to the server periodically using another Redis queue.

## Requirements

- Docker
- Docker Compose
- Redis

## Setup

1. Clone the repository:

   ```bash
   git clone <repository-url>


2. Execute and build by scaling worker to 3 containers

   ```bash
   docker-compose up --scale worker=3 --build

## Refer to manager executor and worker documentation

[Link to manager README](./manager/README.md)


[Link to worker README](./worker/README.md)


## Output log

This log shows how messages are published and instantly consumed by available workers and how each worker reports its tasks status back to the manager

```bash
Attaching to go-executor_redis_1, go-executor_manager_1, go-executor_worker_1, go-executor_worker_3, go-executor_worker_2

manager_1  | Connected to Redis: PONG
manager_1  | Published message: Message 1
manager_1  | Published message: Message 2
manager_1  | Published message: Message 3
manager_1  | Published message: Message 4
manager_1  | Published message: Message 5


worker_1   | Connected to Redis: PONG
worker_1   | Worker consumed message: Message 1
worker_1   | Reported task status: Message 1 - Completed 50% Task
manager_1  | Manager consumed message regarding worker status on Task: Message 1: Completed 50% Task
worker_1   | Reported task status: Message 1 - Completed 100% Task
worker_1   | Worker finished performing Task: Message 1
manager_1  | Manager consumed message regarding worker status on Task: Message 1: Completed 100% Task

worker_1   | Worker consumed message: Message 4
worker_1   | Reported task status: Message 4 - Completed 50% Task
manager_1  | Manager consumed message regarding worker status on Task: Message 4: Completed 50% Task
worker_1   | Reported task status: Message 4 - Completed 100% Task
worker_1   | Worker finished performing Task: Message 4
manager_1  | Manager consumed message regarding worker status on Task: Message 4: Completed 100% Task


worker_2   | Connected to Redis: PONG
worker_2   | Worker consumed message: Message 3
worker_2   | Reported task status: Message 3 - Completed 50% Task
manager_1  | Manager consumed message regarding worker status on Task: Message 3: Completed 50% Task
worker_2   | Reported task status: Message 3 - Completed 100% Task
worker_2   | Worker finished performing Task: Message 3



worker_3   | Connected to Redis: PONG
worker_3   | Worker consumed message: Message 2
worker_3   | Reported task status: Message 2 - Completed 50% Task
manager_1  | Manager consumed message regarding worker status on Task: Message 2: Completed 50% Task
worker_3   | Reported task status: Message 2 - Completed 100% Task
worker_3   | Worker finished performing Task: Message 2
manager_1  | Manager consumed message regarding worker status on Task: Message 2: Completed 100% Task


worker_3   | Worker consumed message: Message 5
worker_3   | Reported task status: Message 5 - Completed 50% Task
manager_1  | Manager consumed message regarding worker status on Task: Message 5: Completed 50% Task
worker_3   | Reported task status: Message 5 - Completed 100% Task
worker_3   | Worker finished performing Task: Message 5
manager_1  | Manager consumed message regarding worker status on Task: Message 5: Completed 100% Task

