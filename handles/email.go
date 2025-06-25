package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sahilrana27582/go-task-queue/models/payloads"
)

func HandleEmailTask(p payloads.EmailPayload) error {
	log.Printf("📧 [EmailHandler] Sending email to: %s | Subject: %s", p.To, p.Subject)

	if !strings.Contains(p.To, "@") {
		return fmt.Errorf("invalid email address: %s", p.To)
	}
	if p.Subject == "" || p.Body == "" {
		return fmt.Errorf("email subject/body cannot be empty")
	}

	time.Sleep(1 * time.Second)
	log.Printf("✅ [EmailHandler] Email successfully sent to %s", p.To)
	return nil
}
