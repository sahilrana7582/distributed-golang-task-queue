package models

import (
	"encoding/json"
	"time"
)

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
	StatusCancelled  Status = "CANCELLED"
	StatusUnknown    Status = "UNKNOWN"
	StatusRunning    Status = "RUNNING"
)

type Task struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Payload    json.RawMessage `json:"payload"`
	Status     Status          `json:"status"`
	Attempts   int             `json:"attempts"`
	MaxRetries int             `json:"max_retries"`
	CreatedAt  time.Time       `json:"created_at"`
}

func NewTask(id, taskType string, payload any) (*Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	t := &Task{
		ID:         id,
		Type:       taskType,
		Payload:    json.RawMessage(data),
		Status:     StatusPending,
		Attempts:   0,
		MaxRetries: 3,
		CreatedAt:  time.Now(),
	}

	return t, nil
}
