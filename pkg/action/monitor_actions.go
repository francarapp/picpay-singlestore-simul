package action

import (
	"sync"

	"github.com/francarapp/picpay-singlestore-simul/pkg/repo"
)

var Monitor = struct {
	Creations int
	Batchs    int
	AvgTime   int64
	MinTime   int64
	MaxTime   int64
}{0, 0, 0, 999999999, 0}

var mutex sync.Mutex

func done(exec repo.MethodExec, qtd int, millis int64) {
	mutex.Lock()
	defer mutex.Unlock()
	switch exec {
	case repo.CreateExec:
		Monitor.Creations += qtd
		Monitor.Batchs++
		if millis > Monitor.MaxTime {
			Monitor.MaxTime = millis
		}
		if millis < Monitor.MinTime {
			Monitor.MinTime = millis
		}
		Monitor.AvgTime = Monitor.AvgTime + int64((millis-Monitor.AvgTime)/int64(Monitor.Batchs))
	}
}
