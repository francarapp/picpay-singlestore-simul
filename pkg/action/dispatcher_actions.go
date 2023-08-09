package action

import (
	"context"

	"github.com/francarapp/picpay-singlestore-simul/pkg/repo"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
 * Dispatcher configuration.
 */
var config *DispatchConfig
var dispatchChannel chan Action
var repoBuffer []repo.EventRepo

type DispatchConfig struct {
	ChanSize    int
	ThreadsSize int
	BatchSize   int
	DB          *gorm.DB
}

func InitDispatching(cfg *DispatchConfig) error {
	config = cfg
	dispatchChannel = make(chan Action, cfg.ChanSize)
	repoBuffer = []repo.EventRepo{}
	for i := 0; i < cfg.ThreadsSize; i++ {
		repoBuffer = append(repoBuffer, repo.NewGormEventRepo(config.DB.Session(&gorm.Session{
			PrepareStmt:     true,
			CreateBatchSize: config.BatchSize,
			SkipHooks:       true,
			Logger:          logger.Default.LogMode(logger.Silent),
		}), config.BatchSize, done))
	}

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
	dispatchChannel <- action
	return nil
}

func Flush(ctx context.Context) error {
	for _, repo := range repoBuffer {
		repo.Flush(ctx)
	}
	return nil
}
