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
	Create(context.Context, *domain.Event) error
	ForceFlush(ctx context.Context) error
	Flush(ctx context.Context) error

	QueryRTCount(ctx context.Context, events []string, start, end string) error
	QueryRTSum(ctx context.Context, events []string, start, end string) error
	QueryRTLabels(ctx context.Context, events []string, start, end string, labels []string) error
}

type FAfterExec func(method MethodExec, idx int, qtd int, millis int64)
type MethodExec string

var (
	CreateExec        = MethodExec("Create")
	QueryRTCountExec  = MethodExec("QueryRTCount")
	QueryRTSumExec    = MethodExec("QueryRTSum")
	QueryRTLabelsExec = MethodExec("QueryRTExec")
)

func NewGormEventRepo(idx int, db *gorm.DB, batchSize int, after FAfterExec) EventRepo {
	return &gormEventRepo{
		Index:     idx,
		DB:        db,
		BatchSize: batchSize,
		FAfter:    after,
	}
}

type gormEventRepo struct {
	MutexCreate sync.Mutex
	MutexFlush  sync.Mutex
	Index       int
	DB          *gorm.DB
	BatchSize   int
	Buffer      []*domain.Event
	FAfter      FAfterExec
}

func (repo *gormEventRepo) Create(ctx context.Context, event *domain.Event) error {
	repo.MutexCreate.Lock()
	defer repo.MutexCreate.Unlock()
	repo.Buffer = append(repo.Buffer, event)
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
	start := time.Now()
	tx := repo.DB.CreateInBatches(repo.Buffer, repo.BatchSize)
	repo.FAfter(CreateExec, repo.Index, len(repo.Buffer), time.Since(start).Milliseconds())
	repo.Buffer = []*domain.Event{}
	return tx.Error
}

func (repo *gormEventRepo) QueryRTCount(ctx context.Context, events []string, start, end string) error {
	timestamp := time.Now()

	var tx *gorm.DB
	results := make([]map[string]interface{}, 0)
	tx = repo.DB.Exec(
		RTSelectCount,
		start, end, events,
	)
	if tx.Error != nil {
		fmt.Printf("Failed: %s", tx.Error)
	} else {
		fmt.Printf("Payload: %s", results)
	}
	repo.FAfter(QueryRTCountExec, repo.Index, 1, time.Since(timestamp).Milliseconds())

	return nil
}

func (repo *gormEventRepo) QueryRTSum(ctx context.Context, events []string, start, end string) error {
	timestamp := time.Now()

	var tx *gorm.DB
	results := make([]map[string]interface{}, 0)
	tx = repo.DB.Raw(
		RTSelectSum,
		start, end, events,
	).Find(&results)
	if tx.Error != nil {
		fmt.Printf("Failed: %s", tx.Error)
	} else {
		fmt.Printf("Payload: %s", results)
	}
	repo.FAfter(QueryRTCountExec, repo.Index, 1, time.Since(timestamp).Milliseconds())

	return nil
}

func (repo *gormEventRepo) QueryRTLabels(ctx context.Context, events []string, start, end string, labels []string) error {
	return nil
}

var (
	RTSelectCount = `
	select event_name, dt_created_min, format(count(*), 0)
	from n_event 
	where 
	  dt_created_min between ? and ?
	  and  event_name in (?)
	group by event_name, dt_created_min
	`
	RTSelectSum = `
	select event_name, dt_created_min, format(sum(payload::$valor), 0)
	from n_event 
	where 
	  dt_created_min between ? and ?
	  and  event_name in (?)
	group by event_name, dt_created_min
	`
)
