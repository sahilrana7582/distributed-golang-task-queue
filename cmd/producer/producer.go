package producer

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/sahilrana27582/go-task-queue/internal/queue"
	"github.com/sahilrana27582/go-task-queue/models"
	"github.com/sahilrana27582/go-task-queue/models/payloads"
)

func ProducePayload(redisQueue *queue.RedisQueue, ctx context.Context) {

	file, err := os.OpenFile("./cmd/producer/producer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger := log.New(file, "Producer: ", log.LstdFlags|log.Lshortfile)

	go func() {
		defer file.Close()

		logger.Printf("====================     %s    =======================", time.Now().Format("2006-01-02 15:04:05"))
		for {
			select {
			case <-ctx.Done():
				logger.Println("❌ Producer stopped")
				logger.Println("===========================================")
				return
			default:
				taskType, payloadData := randomPayload()
				task, err := models.NewTask("task-"+time.Now().Format("20060102150405"), taskType, payloadData)
				if err != nil {
					logger.Printf("❌ Error creating task: %v", err)
					logger.Println()
					continue
				}

				if err := redisQueue.Enqueue(task); err != nil {
					logger.Printf("❌ Error enqueueing task: %v", err)
					logger.Println()
					continue
				}

				logger.Printf("✅ Enqueued task: ID=%s | Type=%s | Payload=%s", task.ID, task.Type, string(task.Payload)+"\n")
				time.Sleep(5 * time.Second)
			}
		}

	}()
}

func randomPayload() (string, []byte) {
	switch rand.Intn(4) {
	case 0:
		p := payloads.EmailPayload{To: "a@example.com", Subject: "Hello", Body: "Body text"}
		b, _ := json.Marshal(p)
		return "email", b
	case 1:
		p := payloads.ImageResizePayload{ImageURL: "http://img.com/1.jpg", Width: 800, Height: 600}
		b, _ := json.Marshal(p)
		return "image_resize", b
	case 2:
		p := payloads.NotificationPayload{UserID: "user123", Message: "Ping!", Channel: "push", Timestamp: time.Now().Unix()}
		b, _ := json.Marshal(p)
		return "notification", b
	default:
		p := payloads.DataBackupPayload{ResourceID: "res01", BackupType: "full", Destination: "s3://bucket/backup", RequestedBy: "admin"}
		b, _ := json.Marshal(p)
		return "data_backup", b
	}
}
