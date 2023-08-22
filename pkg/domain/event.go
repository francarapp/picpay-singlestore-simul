package domain

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Event struct {
	EventName     string
	EventID       string
	CorrelationID string
	UserID        string
	DtCreated     string
	DtReceived    string
	DtIngested    string
	Labels        string
	Payload       string
}

// TableName overrides the table name used by User to `profiles`
func (Event) TableName() string {
	return "nn_event"
}

func (event *Event) Clone() *Event {
	return &Event{
		EventName:     event.EventName,
		EventID:       uuid.NewString(),
		CorrelationID: event.CorrelationID,
		UserID:        event.UserID,
		DtCreated:     event.DtCreated,
		DtReceived:    event.DtReceived,
		DtIngested:    event.DtIngested,
		Labels:        event.Labels,
		Payload:       event.Payload,
	}
}

func NewSEvent(event *Event, vector []float64) *SEvent {
	sparse := false
	if len(vector) > 0 {
		sparse = true
	}
	return &SEvent{
		Event:  *event,
		Sparse: sparse,
		Vector: Vector{vector},
	}
}

type SEvent struct {
	Event
	Sparse bool `gorm:"-"`
	Vector Vector
}

// TableName overrides the table name used by User to `profiles`
func (SEvent) TableName() string {
	return "s_event"
}

type Vector struct {
	data []float64
}

func (vector Vector) GormDataType() string {
	return "vector"
}

func (vector Vector) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := json.Marshal(vector.data)
	return clause.Expr{
		SQL:  "json_array_pack(?)",
		Vars: []interface{}{string(v)},
	}
}

// Scan implements the sql.Scanner interface
func (vector *Vector) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	return nil
}
