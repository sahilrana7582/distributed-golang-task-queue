package main

import (
	"time"

	"github.com/sahilrana27582/go-task-queue/internal/queue"
)

func main() {
	redisQueue := queue.NewRedisQueue("localhost:6379")

	status, err := redisQueue.Ping()

	if err != nil {
		panic(err)
	}
	println("Redis connection status:", status)

	time.Sleep(2 * time.Second)
	// Example of enqueuing a task

}
