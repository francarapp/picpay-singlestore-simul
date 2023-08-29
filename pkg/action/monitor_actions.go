package action

import (
	"fmt"
	"sync"

	"github.com/francarapp/picpay-singlestore-simul/pkg/repo"
)

type AccumEnum string

var (
	AcNone       = AccumEnum("none")
	AcDispatches = AccumEnum("dispatches")
	AcExecutions = AccumEnum("creations")
	AcBatches    = AccumEnum("batches")
	AcAvg        = AccumEnum("avg")
	AcMin        = AccumEnum("min")
	AcMax        = AccumEnum("max")
)

type AccumulatorHolder interface {
	Add(ac AccumEnum, c int64)
	Get(AccumEnum) int64
	Clean() AccumulatorHolder
	String() string
}

type Accumulator struct {
	Lock   sync.RWMutex
	Holder AccumulatorHolder
}

func NewAccumulator(holder AccumulatorHolder) *Accumulator {
	return &Accumulator{
		Holder: holder,
	}
}

func (accum *Accumulator) Clean() *Accumulator {
	accum.Lock.Lock()
	defer accum.Lock.Unlock()
	accum.Holder.Clean()
	return accum
}

func (accum *Accumulator) String() string {
	return accum.Holder.String()
}

type acCommand struct {
	Ac AccumEnum
	C  int64
}

func (accum *Accumulator) Add(acs ...acCommand) {
	accum.Lock.Lock()
	defer accum.Lock.Unlock()
	if len(acs) == 0 {
		accum.Holder.Add(AcNone, 1)
	}
	for _, comm := range acs {
		accum.Holder.Add(comm.Ac, comm.C)
	}
}

func (accum *Accumulator) Get(ac AccumEnum) int64 {
	accum.Lock.RLock()
	defer accum.Lock.RUnlock()
	return accum.Holder.Get(ac)
}

type IntAccumHolder struct {
	c int64
}

func (holder *IntAccumHolder) Add(ac AccumEnum, c int64) {
	holder.c += c
}

func (holder *IntAccumHolder) Get(AccumEnum) int64 {
	return holder.c
}

func (holder *IntAccumHolder) Clean() AccumulatorHolder {
	holder.c = 0
	return holder
}

func (holder *IntAccumHolder) String() string {
	return holder.String()
}

var MonitorActionDispatch = NewAccumulator(&IntAccumHolder{})
var MonitorActionCreate = NewAccumulator(&IntAccumHolder{})
var MonitorRepoCreate = NewAccumulator((&BatchAccumHolder{}).Clean())
var MonitorRepoQuery = NewAccumulator((&BatchAccumHolder{}).Clean())

type BatchAccumHolder struct {
	Executions int64
	Batchs     int64
	AvgTime    int64
	MinTime    int64
	MaxTime    int64
}

func (holder *BatchAccumHolder) Add(ac AccumEnum, c int64) {
	switch ac {
	case AcExecutions:
		holder.Executions += c
	case AcBatches:
		holder.Batchs += 1
	case AcAvg:
		holder.AvgTime = c
	case AcMin:
		holder.MinTime = c
	case AcMax:
		holder.MaxTime = c
	}
}

func (holder *BatchAccumHolder) Get(ac AccumEnum) int64 {
	switch ac {
	case AcExecutions:
		return holder.Executions
	case AcBatches:
		return holder.Batchs
	case AcAvg:
		return holder.AvgTime
	case AcMin:
		return holder.MinTime
	case AcMax:
		return holder.MaxTime
	}
	return 0
}

func (holder *BatchAccumHolder) Clean() AccumulatorHolder {
	holder.Executions = 0
	holder.Batchs = 0
	holder.AvgTime = 0
	holder.MinTime = 9999999
	holder.MaxTime = 0
	return holder
}

func (holder *BatchAccumHolder) String() string {
	return fmt.Sprintf("Execs: %d, Batches: %d, Avg: %d", holder.Executions, holder.Batchs, holder.AvgTime)
}

func done(exec repo.MethodExec, idx int, qtd int, millis int64) {
	switch exec {
	case repo.CreateExec:
		batches := MonitorRepoCreate.Get(AcBatches)
		avg := MonitorRepoCreate.Get(AcAvg)
		avg = avg + int64((millis-avg)/(batches+1))
		MonitorRepoCreate.Add(
			acCommand{AcExecutions, int64(qtd)},
			acCommand{AcBatches, 1},
			acCommand{AcMax, IfThenElse(millis > MonitorRepoCreate.Get(AcMax), millis, MonitorRepoCreate.Get(AcMax))},
			acCommand{AcMin, IfThenElse(millis < MonitorRepoCreate.Get(AcMin), millis, MonitorRepoCreate.Get(AcMin))},
			acCommand{AcAvg, avg},
		)
	case repo.QueryRTCountExec, repo.QueryRTSumExec, repo.QueryRTSumValueExec, repo.QueryRTLabelsExec:
		fmt.Printf("%s[%d]: %d millis \n", exec, idx, millis)
		executions := MonitorRepoQuery.Get(AcExecutions)
		avg := MonitorRepoQuery.Get(AcAvg)
		avg = avg + int64((millis-avg)/(executions+1))
		MonitorRepoQuery.Add(
			acCommand{AcExecutions, int64(qtd)},
			acCommand{AcBatches, 0},
			acCommand{AcMax, IfThenElse(millis > MonitorRepoQuery.Get(AcMax), millis, MonitorRepoQuery.Get(AcMax))},
			acCommand{AcMin, IfThenElse(millis < MonitorRepoQuery.Get(AcMin), millis, MonitorRepoQuery.Get(AcMin))},
			acCommand{AcAvg, avg},
		)
	}
}

func Clean() {
	MonitorActionCreate.Clean()
	MonitorActionDispatch.Clean()
	MonitorRepoCreate.Clean()
	MonitorRepoQuery.Clean()
}

func IfThenElse(condition bool, a int64, b int64) int64 {
	if condition {
		return a
	}
	return b
}
