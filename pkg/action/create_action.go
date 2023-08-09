package action

import (
	"context"
	"time"

	"github.com/francarapp/picpay-singlestore-simul/pkg/domain"
	"github.com/francarapp/picpay-singlestore-simul/pkg/repo"
)

type Action interface {
	Do(ctx context.Context) error
}

var repoIndex = 0

func Create(event *domain.Event) Action {
	repoIndex = (repoIndex + 1) % config.ThreadsSize
	return &eventCreateAct{
		Event:     event,
		Repo:      repoBuffer[repoIndex],
		BatchSize: config.BatchSize,
	}
}

func Query(parameters map[string]string) Action {
	return nil
}

func QueryBetween(dtini, dtend time.Time, parameters map[string]string) Action {
	return nil
}

// =================================================
// ===============   CREATE ACTION   ===============

type eventCreateAct struct {
	Event     *domain.Event
	Repo      repo.EventRepo
	BatchSize int
}

func (act *eventCreateAct) Do(ctx context.Context) error {
	MonitorCreate.Add(1)
	return act.Repo.Create(ctx, act.Event)
}

// ================================================
// ===============   QUERY ACTION   ===============
