package action

import (
	"context"
	"sync"

	"github.com/francarapp/picpay-singlestore-simul/pkg/repo"
	"gorm.io/gorm"
)

/*
 * Dispatcher configuration.
 */
var config *DispatchConfig
var dispatchChannel chan Action
var repoBuffer *RepoBuffer

type DispatchConfig struct {
	ChanSize    int
	ThreadsSize int
	BatchSize   int
	DB          *gorm.DB
}

func InitDispatching(cfg *DispatchConfig) error {
	config = cfg
	dispatchChannel = make(chan Action, cfg.ChanSize)
	repoBuffer = newRepoBuffer(cfg)

	for i := 0; i < cfg.ThreadsSize; i++ {
		go func() {
			for act := range dispatchChannel {
				ctx := context.Background()
				act.Do(ctx)
			}
		}()
	}
	return nil
}

func Dispatch(action Action) error {
	MonitorActionDispatch.Add()
	dispatchChannel <- action
	return nil
}

func ForceFlush(ctx context.Context) error {
	repoBuffer.ForceFlush(ctx)
	return nil
}

type RepoBuffer struct {
	Lock   sync.RWMutex
	Index  int
	Buffer []repo.EventRepo
}

func newRepoBuffer(cfg *DispatchConfig) *RepoBuffer {
	repoBuffer := RepoBuffer{Buffer: []repo.EventRepo{}}
	for i := 0; i < cfg.ThreadsSize; i++ {
		repoBuffer.Buffer = append(repoBuffer.Buffer, repo.NewGormEventRepo(i, config.DB.Session(&gorm.Session{
			PrepareStmt:     true,
			CreateBatchSize: config.BatchSize,
			SkipHooks:       true,
		}), config.BatchSize, done))
	}
	return &repoBuffer
}

func (repob *RepoBuffer) Next() repo.EventRepo {
	repob.Lock.Lock()
	defer repob.Lock.Unlock()
	repob.Index = (repob.Index + 1) % len(repob.Buffer)
	return repob.Buffer[repob.Index]
}

func (repob *RepoBuffer) ForceFlush(ctx context.Context) {
	repob.Lock.Lock()
	defer repob.Lock.Unlock()
	for _, repo := range repob.Buffer {
		repo.ForceFlush(ctx)
	}
}
