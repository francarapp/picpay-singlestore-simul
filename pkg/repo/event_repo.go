package repo

import (
	"context"
	"sync"
	"time"

	"github.com/francarapp/picpay-singlestore-simul/pkg/domain"
	"gorm.io/gorm"
)

type EventRepo interface {
	Create(context.Context, *domain.Event) error
	Flush(ctx context.Context) error
}

type FAfterExec func(method MethodExec, qtd int, millis int64)
type MethodExec string

var (
	CreateExec = MethodExec("Create")
	QueryExec  = MethodExec("Query")
)

func NewGormEventRepo(db *gorm.DB, batchSize int, after FAfterExec) EventRepo {
	return &gormEventRepo{
		DB:        db,
		BatchSize: batchSize,
		After:     after,
	}
}

type gormEventRepo struct {
	Mutex     sync.Mutex
	DB        *gorm.DB
	BatchSize int
	Buffer    []*domain.Event
	After     FAfterExec
}

func (repo *gormEventRepo) Create(ctx context.Context, event *domain.Event) error {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()
	repo.Buffer = append(repo.Buffer, event)
	if len(repo.Buffer) >= repo.BatchSize {
		return repo.Flush(ctx)
	}
	return nil
}

func (repo *gormEventRepo) Flush(ctx context.Context) error {
	start := time.Now()
	tx := repo.DB.CreateInBatches(repo.Buffer, repo.BatchSize)
	repo.After(CreateExec, len(repo.Buffer), time.Since(start).Milliseconds())
	repo.Buffer = []*domain.Event{}
	return tx.Error
}
