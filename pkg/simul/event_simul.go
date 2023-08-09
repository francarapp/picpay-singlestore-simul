package simul

import (
	"context"
	"time"

	"github.com/francarapp/picpay-singlestore-simul/pkg/domain"
	"github.com/google/uuid"
)

func NewEvent(ctx context.Context) *domain.Event {
	correlationID := CorrelationID("999")
	if ctx.Value(CorrelationKey) != nil {
		correlationID = ctx.Value(CorrelationKey).(CorrelationID)
	}
	eventName := genEventName()
	ev := domain.Event{
		EventName:     eventName,
		EventID:       uuid.New().String(),
		CorrelationID: string(correlationID),
		DtCreated:     time.Now(),
		// DtReceived:    time.Now(),
		// DtIngested:    time.Now(),
		Labels:  genLabels(eventName),
		Payload: genPayload(eventName),
	}
	return &ev
}
