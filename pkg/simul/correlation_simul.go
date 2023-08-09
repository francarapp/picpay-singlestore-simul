package simul

import "context"

var (
	UserKey        = UserID("UserID")
	CorrelationKey = CorrelationID("CorrelationID")
)

type CorrelationID string
type UserID string

func CorrelateContext(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, CorrelationKey, CorrelationID(correlationId))
}

func UserContext(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, UserKey, UserID(userId))
}
