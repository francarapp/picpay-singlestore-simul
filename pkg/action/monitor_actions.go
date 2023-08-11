package action

import (
	"sync"

	"github.com/francarapp/picpay-singlestore-simul/pkg/repo"
)

type Accumulator struct {
	Lock  sync.RWMutex
	Count int
}

func (accum *Accumulator) Clean() {
	accum.Lock.Lock()
	defer accum.Lock.Unlock()
	accum.Count = 0
}

func (accum *Accumulator) Add(c int) {
	accum.Lock.Lock()
	defer accum.Lock.Unlock()
	accum.Count += c
}

func (accum *Accumulator) Get() int {
	accum.Lock.RLock()
	defer accum.Lock.RUnlock()
	return accum.Count
}

var MonitorDispatch = &Accumulator{}
var MonitorCreate = &Accumulator{}

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

func Clean() {
	mutex.Lock()
	defer mutex.Unlock()

	Monitor.AvgTime = 0
	Monitor.Batchs = 0
	Monitor.Creations = 0
	Monitor.MaxTime = 0
	Monitor.MinTime = 0

	MonitorCreate.Clean()
	MonitorDispatch.Clean()
}
