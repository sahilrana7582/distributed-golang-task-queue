package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sahilrana27582/go-task-queue/models/payloads"
)

func HandleDataBackupTask(p payloads.DataBackupPayload) error {
	log.Printf("💾 [BackupHandler] Starting backup for resource %s (%s) to %s",
		p.ResourceID, p.BackupType, p.Destination)

	if p.ResourceID == "" || p.BackupType == "" || p.Destination == "" {
		return fmt.Errorf("missing fields in data backup payload")
	}
	if !strings.HasPrefix(p.Destination, "s3://") {
		return fmt.Errorf("unsupported backup destination: %s", p.Destination)
	}

	time.Sleep(3 * time.Second)
	log.Printf("✅ [BackupHandler] Backup completed for resource %s", p.ResourceID)
	return nil
}
