package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/sahilrana27582/go-task-queue/models/payloads"
)

func HandleNotificationTask(p payloads.NotificationPayload) error {
	log.Printf("📣 [NotificationHandler] Sending notification to user %s on channel %s", p.UserID, p.Channel)

	if p.UserID == "" || p.Message == "" || p.Channel == "" {
		return fmt.Errorf("missing fields in notification payload")
	}

	time.Sleep(1 * time.Second)
	log.Printf("✅ [NotificationHandler] Notification sent to %s: %s", p.UserID, p.Message)
	return nil
}
