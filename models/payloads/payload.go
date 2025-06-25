package payloads

import (
	"encoding/json"
	"fmt"
)

type Payload interface {
	Decode(data []byte) (Payload, error)
}

// ------------------- Payload Structs -------------------

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type ImageResizePayload struct {
	ImageURL string `json:"image_url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type NotificationPayload struct {
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	Channel   string `json:"channel"`
	Timestamp int64  `json:"timestamp"`
}

type DataBackupPayload struct {
	ResourceID  string `json:"resource_id"`
	BackupType  string `json:"backup_type"` // e.g., "full", "incremental"
	Destination string `json:"destination"` // e.g., S3 bucket URL
	RequestedBy string `json:"requested_by"`
}

// ------------------- Decode Implementations -------------------

func (EmailPayload) Decode(data []byte) (Payload, error) {
	var p EmailPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("error decoding EmailPayload: %w", err)
	}
	return p, nil
}

func (ImageResizePayload) Decode(data []byte) (Payload, error) {
	var p ImageResizePayload
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("error decoding ImageResizePayload: %w", err)
	}
	return p, nil
}

func (NotificationPayload) Decode(data []byte) (Payload, error) {
	var p NotificationPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("error decoding NotificationPayload: %w", err)
	}
	return p, nil
}

func (DataBackupPayload) Decode(data []byte) (Payload, error) {
	var p DataBackupPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("error decoding DataBackupPayload: %w", err)
	}
	return p, nil
}

// ------------------- Payload Registry -------------------

var payloadRegistry = map[string]Payload{
	"email":        EmailPayload{},
	"image_resize": ImageResizePayload{},
	"notification": NotificationPayload{},
	"data_backup":  DataBackupPayload{},
}

func DecodePayload(taskType string, raw json.RawMessage) (Payload, error) {
	proto, ok := payloadRegistry[taskType]
	if !ok {
		return nil, fmt.Errorf("unsupported task type: %s", taskType)
	}
	return proto.Decode(raw)
}
