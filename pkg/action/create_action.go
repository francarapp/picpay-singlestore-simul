package action

import (
	"context"

	"github.com/francarapp/picpay-singlestore-simul/pkg/domain"
	"github.com/francarapp/picpay-singlestore-simul/pkg/repo"
)

type Action interface {
	Do(ctx context.Context) error
}

func Create(event *domain.Event, vector ...float64) Action {
	sparse := false
	if len(vector) > 0 {
		sparse = true
	}
	return &eventCreateAct{
		Event:     event,
		Sparse:    sparse,
		Vector:    vector,
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

func QueryRTSumValue(events []string, start, end string) Action {
	return &eventQueryAct{
		Repo:   repoBuffer.Next(),
		Exec:   repo.QueryRTSumValueExec,
		Events: events,
		Start:  start,
		End:    end,
	}
}

func QueryMRTCount(events []string, start, end string) Action {
	return &eventQueryAct{
		Repo:   repoBuffer.Next(),
		Exec:   repo.QueryMRTCountExec,
		Events: events,
		Start:  start,
		End:    end,
	}
}

func QueryMRTSum(events []string, start, end string) Action {
	return &eventQueryAct{
		Repo:   repoBuffer.Next(),
		Exec:   repo.QueryMRTSumExec,
		Events: events,
		Start:  start,
		End:    end,
	}
}

// =================================================
// ===============   CREATE ACTION   ===============

type eventCreateAct struct {
	Event     *domain.Event
	Vector    []float64
	Repo      repo.EventRepo
	Sparse    bool
	BatchSize int
}

func (act *eventCreateAct) Do(ctx context.Context) error {
	MonitorActionCreate.Add()
	if act.Sparse {
		return act.Repo.Create(ctx, act.Event, act.Vector...)
	}
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
	case repo.QueryRTSumValueExec:
		return act.Repo.QueryRTSumValue(ctx, act.Events, act.Start, act.End)
	case repo.QueryMRTCountExec:
		return act.Repo.QueryMRTCount(ctx, act.Events, act.Start, act.End)
	case repo.QueryMRTSumExec:
		return act.Repo.QueryMRTSum(ctx, act.Events, act.Start, act.End)
	}
	return nil
}
