package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/sahilrana27582/go-task-queue/models/payloads"
)

func HandleImageResizeTask(p payloads.ImageResizePayload) error {
	log.Printf("🖼️ [ImageResizeHandler] Resizing image: %s to %dx%d", p.ImageURL, p.Width, p.Height)

	if p.ImageURL == "" || p.Width <= 0 || p.Height <= 0 {
		return fmt.Errorf("invalid image resize parameters")
	}

	time.Sleep(2 * time.Second)
	log.Printf("✅ [ImageResizeHandler] Image resized successfully: %s", p.ImageURL)
	return nil
}
