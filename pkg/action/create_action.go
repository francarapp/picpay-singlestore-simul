package action

import (
	"context"

	"github.com/francarapp/picpay-singlestore-simul/pkg/domain"
	"github.com/francarapp/picpay-singlestore-simul/pkg/repo"
)

type Action interface {
	Do(ctx context.Context) error
}

func Create(event *domain.Event) Action {
	return &eventCreateAct{
		Event:     event,
		Repo:      repoBuffer.Next(),
		BatchSize: config.BatchSize,
	}
}

func QueryRTCount(events []string, start, end string) Action {
	return &eventQueryAct{
		Repo:   repoBuffer.Next(),
		Exec:   repo.QueryRTCountExec,
		Events: events,
		Start:  start,
		End:    end,
	}
}

func QueryRTSum(events []string, start, end string) Action {
	return &eventQueryAct{
		Repo:   repoBuffer.Next(),
		Exec:   repo.QueryRTSumExec,
		Events: events,
		Start:  start,
		End:    end,
	}
}

// =================================================
// ===============   CREATE ACTION   ===============

type eventCreateAct struct {
	Event     *domain.Event
	Repo      repo.EventRepo
	BatchSize int
}

func (act *eventCreateAct) Do(ctx context.Context) error {
	MonitorActionCreate.Add()
	return act.Repo.Create(ctx, act.Event)
}

// ================================================
// ===============   QUERY ACTION   ===============

type eventQueryAct struct {
	Repo   repo.EventRepo
	Exec   repo.MethodExec
	Events []string
	Start  string
	End    string
}

func (act *eventQueryAct) Do(ctx context.Context) error {
	MonitorActionCreate.Add()
	switch act.Exec {
	case repo.QueryRTCountExec:
		return act.Repo.QueryRTCount(ctx, act.Events, act.Start, act.End)
	case repo.QueryRTSumExec:
		return act.Repo.QueryRTSum(ctx, act.Events, act.Start, act.End)
	}
	return nil
}
