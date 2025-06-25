package producer

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/sahilrana27582/go-task-queue/internal/queue"
	"github.com/sahilrana27582/go-task-queue/models"
	"github.com/sahilrana27582/go-task-queue/models/payloads"
)

func ProducePayload(redisQueue *queue.RedisQueue, ctx context.Context, numProducers int) {
	file, err := os.OpenFile("./cmd/producer/producer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open producer log file: %v", err)
	}
	logger := log.New(file, "Producer: ", log.LstdFlags|log.Lshortfile)

	logger.Printf("🚀 Starting %d producer(s) at %s", numProducers, time.Now().Format("2006-01-02 15:04:05"))
	defer file.Close()

	wg := sync.WaitGroup{}
	wg.Add(numProducers)

	for i := 0; i < numProducers; i++ {
		go func(producerID int) {

			defer func() {
				if r := recover(); r != nil {
					logger.Printf("❌ Producer-%d encountered an error: %v", producerID, r)
				}
				wg.Done()
			}()

			logger.Printf("🧵 Producer-%d started", producerID)

			for {
				select {
				case <-ctx.Done():
					wg.Done()
					logger.Printf("❌ Producer-%d stopped", producerID)
					return
				default:
					taskType, payloadStruct := randomPayload()

					task, err := models.NewTask(
						"task-"+time.Now().Format("20060102150405")+"-"+randomID(),
						taskType,
						payloadStruct,
					)
					if err != nil {
						logger.Printf("❌ Producer-%d: Task creation failed: %v", producerID, err)
						continue
					}

					if err := redisQueue.Enqueue(task); err != nil {
						logger.Printf("❌ Producer-%d: Enqueue failed: %v", producerID, err)
						continue
					}

					logger.Printf("✅ Producer-%d: Enqueued Task ID=%s | Type=%s", producerID, task.ID, task.Type)

					time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond) // slight jitter
				}
			}
		}(i + 1)

	}

	wg.Wait()
}

func randomPayload() (string, any) {
	switch rand.Intn(4) {
	case 0:
		return "email", payloads.EmailPayload{
			To: "a@example.com", Subject: "Hello", Body: "Body text",
		}
	case 1:
		return "image_resize", payloads.ImageResizePayload{
			ImageURL: "http://img.com/1.jpg", Width: 800, Height: 600,
		}
	case 2:
		return "notification", payloads.NotificationPayload{
			UserID: "user123", Message: "Ping!", Channel: "push", Timestamp: time.Now().Unix(),
		}
	default:
		return "data_backup", payloads.DataBackupPayload{
			ResourceID: "res01", BackupType: "full", Destination: "s3://bucket/backup", RequestedBy: "admin",
		}
	}
}

func randomID() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}
