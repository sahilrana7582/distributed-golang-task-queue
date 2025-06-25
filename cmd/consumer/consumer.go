package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	handlers "github.com/sahilrana27582/go-task-queue/handles"
	"github.com/sahilrana27582/go-task-queue/internal/queue"
	"github.com/sahilrana27582/go-task-queue/models/payloads"
)

func StartWorker(redisQueue *queue.RedisQueue, ctx context.Context, numWorkers int) {
	file, err := os.OpenFile("./cmd/consumer/consumer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open consumer log file: %v", err)
	}
	logger := log.New(file, "Consumer: ", log.LstdFlags|log.Lshortfile)

	logger.Printf("🚀 Starting %d workers at %s", numWorkers, time.Now().Format("2006-01-02 15:04:05"))

	defer file.Close()
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			logger.Printf("🧵 Worker-%d started", workerID)
			defer func() {
				if r := recover(); r != nil {
					logger.Printf("❌ Worker-%d encountered an error: %v", workerID, r)
				}
				wg.Done()
			}()

			for {
				select {
				case <-ctx.Done():
					logger.Printf("❌ Worker-%d stopped", workerID)
					return
				default:
					task, err := redisQueue.Dequeue(2 * time.Second)
					if err != nil {
						if err.Error() == "redis nil" {
							logger.Printf("🔄 Worker-%d: No task, retrying...", workerID)
							continue
						}
						logger.Printf("❌ Worker-%d: Failed to dequeue: %v", workerID, err)
						continue
					}

					logger.Printf("⏳ Worker-%d: Task ID=%s, Type=%s", workerID, task.ID, task.Type)
					if err := dispatch(task.Type, task.Payload); err != nil {
						logger.Printf("❌ Worker-%d: Failed to handle task %s: %v", workerID, task.ID, err)
					} else {
						logger.Printf("✅ Worker-%d: Task %s done", workerID, task.ID)
					}
				}
			}
		}(i + 1)
	}

	wg.Wait()
}

func dispatch(taskType string, rawPayload []byte) error {
	switch taskType {
	case "email":
		var p payloads.EmailPayload
		if err := json.Unmarshal(rawPayload, &p); err != nil {
			return err
		}
		return handlers.HandleEmailTask(p)
	case "image_resize":
		var p payloads.ImageResizePayload
		if err := json.Unmarshal(rawPayload, &p); err != nil {
			return err
		}
		return handlers.HandleImageResizeTask(p)
	case "notification":
		var p payloads.NotificationPayload
		if err := json.Unmarshal(rawPayload, &p); err != nil {
			return err
		}
		return handlers.HandleNotificationTask(p)
	case "data_backup":
		var p payloads.DataBackupPayload
		if err := json.Unmarshal(rawPayload, &p); err != nil {
			return err
		}
		return handlers.HandleDataBackupTask(p)
	default:
		return fmt.Errorf("unknown task type: %s", taskType)
	}
}
