package repo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/francarapp/picpay-singlestore-simul/pkg/domain"
	"gorm.io/gorm"
)

type EventRepo interface {
	Create(context.Context, *domain.Event, ...float64) error
	ForceFlush(ctx context.Context) error
	Flush(ctx context.Context) error

	QueryRTCount(ctx context.Context, events []string, start, end string) error
	QueryRTSum(ctx context.Context, events []string, start, end string) error
	QueryRTLabels(ctx context.Context, events []string, start, end string, labels []string) error

	QueryMRTCount(ctx context.Context, events []string, start, end string) error
	QueryMRTSum(ctx context.Context, events []string, start, end string) error
	QueryMRTLabels(ctx context.Context, events []string, start, end string, labels []string) error
}

type FAfterExec func(method MethodExec, idx int, qtd int, millis int64)
type MethodExec string

var (
	CreateExec         = MethodExec("Create")
	QueryRTCountExec   = MethodExec("QueryRTCount")
	QueryRTSumExec     = MethodExec("QueryRTSum")
	QueryRTLabelsExec  = MethodExec("QueryRTExec")
	QueryMRTCountExec  = MethodExec("QueryMRTCount")
	QueryMRTSumExec    = MethodExec("QueryMRTSum")
	QueryMRTLabelsExec = MethodExec("QueryMRTExec")
)

func NewGormEventRepo(idx int, db *gorm.DB, sparse bool, batchSize int, after FAfterExec) EventRepo {
	return &gormEventRepo{
		Index:     idx,
		DB:        db,
		Sparse:    sparse,
		BatchSize: batchSize,
		FAfter:    after,
	}
}

type gormEventRepo struct {
	MutexCreate sync.Mutex
	MutexFlush  sync.Mutex
	Index       int
	DB          *gorm.DB
	Sparse      bool
	BatchSize   int
	Buffer      []*domain.SEvent
	FAfter      FAfterExec
}

func (repo *gormEventRepo) Create(ctx context.Context, event *domain.Event, vector ...float64) error {
	repo.MutexCreate.Lock()
	defer repo.MutexCreate.Unlock()
	repo.Buffer = append(repo.Buffer, domain.NewSEvent(event, vector))
	if len(repo.Buffer) >= repo.BatchSize {
		return repo.Flush(ctx)
	}
	return nil
}

func (repo *gormEventRepo) ForceFlush(ctx context.Context) error {
	repo.MutexCreate.Lock()
	defer repo.MutexCreate.Unlock()
	return repo.Flush(ctx)
}

func (repo *gormEventRepo) Flush(ctx context.Context) error {
	repo.MutexFlush.Lock()
	defer repo.MutexFlush.Unlock()

	buffer := repo.castIf()

	start := time.Now()
	tx := repo.DB.CreateInBatches(buffer, repo.BatchSize)
	repo.FAfter(CreateExec, repo.Index, len(repo.Buffer), time.Since(start).Milliseconds())
	repo.Buffer = []*domain.SEvent{}
	return tx.Error
}

func (repo *gormEventRepo) castIf() interface{} {
	if repo.Sparse {
		return repo.Buffer
	}

	nbuffer := []*domain.Event{}
	for _, ev := range repo.Buffer {
		nbuffer = append(nbuffer, &ev.Event)
	}
	return nbuffer
}

func (repo *gormEventRepo) QueryRTCount(ctx context.Context, events []string, start, end string) error {
	repo.MutexCreate.Lock()
	defer repo.MutexCreate.Unlock()
	timestamp := time.Now()

	var tx *gorm.DB
	tx = repo.DB.Exec(
		RTSelectCount,
		start, end, events,
	)
	if tx.Error != nil {
		fmt.Printf("Failed: %s", tx.Error)
	}
	repo.FAfter(QueryRTCountExec, repo.Index, 1, time.Since(timestamp).Milliseconds())

	return nil
}

func (repo *gormEventRepo) QueryRTSum(ctx context.Context, events []string, start, end string) error {
	timestamp := time.Now()

	var tx *gorm.DB
	tx = repo.DB.Exec(
		RTSelectSum,
		start, end, events,
	)
	if tx.Error != nil {
		fmt.Printf("Failed: %s", tx.Error)
	}
	repo.FAfter(QueryRTCountExec, repo.Index, 1, time.Since(timestamp).Milliseconds())

	return nil
}

func (repo *gormEventRepo) QueryRTLabels(ctx context.Context, events []string, start, end string, labels []string) error {
	return nil
}

func (repo *gormEventRepo) QueryMRTCount(ctx context.Context, events []string, start, end string) error {
	repo.MutexCreate.Lock()
	defer repo.MutexCreate.Unlock()
	timestamp := time.Now()

	var tx *gorm.DB
	tx = repo.DB.Exec(
		MRTSelectCount,
		start, end, events,
	)
	if tx.Error != nil {
		fmt.Printf("Failed: %s", tx.Error)
	}
	repo.FAfter(QueryRTCountExec, repo.Index, 1, time.Since(timestamp).Milliseconds())

	return nil
}

func (repo *gormEventRepo) QueryMRTSum(ctx context.Context, events []string, start, end string) error {
	timestamp := time.Now()

	var tx *gorm.DB
	tx = repo.DB.Exec(
		MRTSelectSum,
		start, end, events,
	)
	if tx.Error != nil {
		fmt.Printf("Failed: %s", tx.Error)
	}
	repo.FAfter(QueryRTCountExec, repo.Index, 1, time.Since(timestamp).Milliseconds())

	return nil
}

func (repo *gormEventRepo) QueryMRTLabels(ctx context.Context, events []string, start, end string, labels []string) error {
	return nil
}

var (
	RTSelectCount = `
	select event_name, dt_created_min, format(count(*), 0)
	from nn_event 
	where 
	  dt_created_min between ? and ?
	  and  event_name in (?)
	group by event_name, dt_created_min
	`
	MRTSelectCount = `
	select event_name, dt_created_min, format(count(*), 0)
	from m_event 
	where 
	  dt_created_min between ? and ?
	  and  event_name in (?)
	group by event_name, dt_created_min
	`
	RTSelectSum = `
	select event_name, dt_created_min, format(sum(payload::$valor), 0)
	from nn_event 
	where 
	  dt_created_min between ? and ?
	  and  event_name in (?)
	group by event_name, dt_created_min
	`
	MRTSelectSum = `
	select event_name, dt_created_min, format(sum(payload::$valor), 0)
	from m_event 
	where 
	  dt_created_min between ? and ?
	  and  event_name in (?)
	group by event_name, dt_created_min
	`
)
