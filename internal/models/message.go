package models

import "time"

// ErrorMessage
type ErrorMessage struct {
	MessageID string    `json:"messageId"`
	Offset    int64     `json:"offset"`
	Partition int       `json:"partition"`
	Topic     string    `json:"topic"`
	Error     string    `json:"error"`
	Time      time.Time `json:"time"`
}
