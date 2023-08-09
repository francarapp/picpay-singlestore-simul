package domain

import (
	"time"
)

type Event struct {
	EventName     string
	EventID       string
	CorrelationID string
	UserID        string
	DtCreated     time.Time
	DtReceived    time.Time
	DtIngested    time.Time
	Labels        string
	Payload       string
}
