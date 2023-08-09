package simul

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/francarapp/picpay-singlestore-simul/pkg/domain"
	"github.com/google/uuid"
)

func NewEvent(ctx context.Context) *domain.Event {
	userID := UserID("888")
	if ctx.Value(UserKey) != nil {
		userID = ctx.Value(UserKey).(UserID)
	}
	correlationID := CorrelationID("999")
	if ctx.Value(CorrelationKey) != nil {
		correlationID = ctx.Value(CorrelationKey).(CorrelationID)
	}
	eventName := genEventName()
	ev := domain.Event{
		EventName:     eventName,
		EventID:       uuid.New().String(),
		UserID:        string(userID),
		CorrelationID: string(correlationID),
		DtCreated:     time.Now().Format("2006-01-02 15:04:05.999"),
		DtReceived:    time.Now().Format("2006-01-02 15:04:05.999"),
		DtIngested:    time.Now().Format("2006-01-02 15:04:05.999"),
		Labels:        genLabels(eventName),
		Payload:       genPayload(eventName),
	}
	return &ev
}

var UsersQtd = 10000000

func GenUserID() string {
	return fmt.Sprintf("consumer_%d", rand.Intn(UsersQtd))
}
