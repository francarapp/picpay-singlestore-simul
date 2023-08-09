package simul

import "context"

var CorrelationKey = CorrelationID("CorrelationID")

type CorrelationID string

func CorrelateContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, CorrelationKey, "")
}
