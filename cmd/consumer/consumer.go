package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	handlers "github.com/sahilrana27582/go-task-queue/handles"
	"github.com/sahilrana27582/go-task-queue/internal/queue"
	"github.com/sahilrana27582/go-task-queue/models/payloads"
)

func StartWorker(redisQueue *queue.RedisQueue, ctx context.Context) {
	file, err := os.OpenFile("./cmd/consumer/consumer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open consumer log file: %v", err)
	}
	logger := log.New(file, "Consumer: ", log.LstdFlags|log.Lshortfile)

	go func() {
		logger.Printf("====================     %s    =======================", time.Now().Format("2006-01-02 15:04:05"))
		logger.Println("🚀 Consumer started, waiting for tasks...")
		defer file.Close()
		for {
			select {
			case <-ctx.Done():
				logger.Println("❌ Consumer stopped")
				logger.Println("===========================================")
				return
			default:
				task, err := redisQueue.Dequeue(1 * time.Second)
				if err != nil {
					if err.Error() == "redis nil" {
						logger.Println("🔄 No tasks in queue, retrying...")
						time.Sleep(1 * time.Second)
						continue
					}
					logger.Printf("❌ Failed to dequeue task: %v", err)
					time.Sleep(1 * time.Second)
					continue
				}

				logger.Printf("⏳ Processing task ID=%s, Type=%s, Payload=%s", task.ID, task.Type, task.Payload)
				fmt.Println("⏳ Processing task ID=", task.ID, "Type=", task.Type, "Payload=", string(task.Payload))

				if err := dispatch(task.Type, task.Payload); err != nil {
					logger.Printf("❌ Error handling task %s: %v", task.ID, err)
				} else {
					logger.Printf("✅ Task %s completed successfully", task.ID)
				}
			}
		}
	}()
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
