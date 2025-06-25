package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sahilrana27582/go-task-queue/cmd/consumer"
	"github.com/sahilrana27582/go-task-queue/cmd/producer"
	"github.com/sahilrana27582/go-task-queue/internal/queue"
)

func main() {
	redisQueue := queue.NewRedisQueue("localhost:6379")
	ctx, cancel := context.WithCancel(context.Background())

	const consumerCount = 5
	defer cancel()

	done := make(chan struct{})

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigs
		log.Printf("🛑 Received signal: %s. Initiating shutdown...", sig)

		cancel() // cancel the context, notifying all components
		close(done)
	}()

	go producer.ProducePayload(redisQueue, ctx, 10)
	go consumer.StartWorker(redisQueue, ctx, 100)

	log.Println("🚀 All services started. Waiting for shutdown signal...")

	<-done

	len, err := redisQueue.Length()
	if err != nil {
		log.Fatalf("❌ Error getting queue length: %v", err)
	}

	fmt.Println("Length of task queue:", len)

	time.Sleep(2 * time.Second)

	log.Println("✅ Graceful shutdown complete.")
}
