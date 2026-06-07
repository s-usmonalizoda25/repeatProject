package models

import "time"

type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	UserID    int       `json:"user_id"`
}
