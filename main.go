package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/francarapp/picpay-singlestore-simul/pkg/action"
	"github.com/francarapp/picpay-singlestore-simul/pkg/simul"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func _main() {
	vector := []float64{0.1, 0.2, 0.3}
	v, _ := json.Marshal(vector)
	fmt.Printf("JSON %s\n", string(v))
}

func main() {
	var codFlag = flag.String("cod", "0", "simul code")
	var threadsFlag = flag.Int("threads", 4, "Paralel instances")
	var createFlag = flag.Bool("create", true, "Create events")
	var deamonFlag = flag.Bool("deamon", false, "Continuous run")
	var silentFlag = flag.Bool("silent", true, "Silent queries")

	var createSparseFlag = flag.Bool("createSparse", false, "generate sparse vector")
	var createQtdFlag = flag.Int("createQtd", 100, "Qtd of create events")
	var createBatchFlag = flag.Int("createBatch", 10, "Qtd batch")
	var createCorrelationsFlag = flag.Int("createCorrelation", 100, "Events with same correlation")

	var queryQtdFlag = flag.Int("queryQtd", 10, "Qtd of events")
	var queryEventsFlag = flag.Int("queryEvents", 5, "Qtd of events")
	var querySelectFlag = flag.String("querySelect", "RTCount", "Query Code")
	var queryStartFlag = flag.String("queryStart", "2023-08-11 12:00:00", "Start date/time")
	var queryEndFlag = flag.String("queryEnd", "2023-08-13 19:00:00", "End data/time")

	flag.Parse()

	config()

	alogger := logger.Default.LogMode(logger.Silent)
	if !*silentFlag {
		alogger = logger.Default.LogMode(logger.Info)
	}
	user := "ingest_events"
	if !*createFlag {
		user = "rt_pp1min"
	}
	db, err := gorm.Open(mysql.Open(
		fmt.Sprintf("%s:12345@tcp(10.164.47.110:3306)/events?parseTime=true&loc=UTC", user),
	), &gorm.Config{
		Logger: alogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		}},
	)
	if err != nil {
		panic("failed to connect database")
	}

	action.InitDispatching(&action.DispatchConfig{
		ChanSize:    1000000,
		ThreadsSize: *threadsFlag,
		BatchSize:   *createBatchFlag,
		Sparse:      *createSparseFlag,
		DB:          db,
	})

	FDo := func() {
		if *createFlag {
			create(db, *createSparseFlag, *codFlag, *threadsFlag, *createQtdFlag, *createBatchFlag, *createCorrelationsFlag)
		} else {
			query(db, *codFlag, *threadsFlag, *queryQtdFlag,
				*queryEventsFlag, *querySelectFlag,
				*queryStartFlag, *queryEndFlag)
		}
	}

	if *deamonFlag {
		for i := 0; ; i++ {
			fmt.Println("\n\n***")
			fmt.Printf("*** SIMUL_%s EXECUTION %d \n", *codFlag, i)
			fmt.Println("***")
			FDo()
		}
	} else {
		FDo()
	}

}

func create(db *gorm.DB, sparse bool, execCod string, threads int, qtdEvs int, batchSize int, maxCorrelations int) error {
	fmt.Printf("*****     CREATING Events[%d] Threads[%d]\n", qtdEvs, threads)
	start := time.Now()
	ctx := context.Background()
	producers := int(threads / 2)
	if producers == 0 {
		producers = 1
	}

	for i := 0; i < producers; i++ {
		count := int(qtdEvs / producers)
		if qtdEvs%producers != 0 && i == 0 {
			count += qtdEvs % producers
		}
		if sparse {
			go produceSparse(ctx, sparse, i, count, maxCorrelations)
		} else {
			go produceCreate(ctx, sparse, i, count, maxCorrelations)
		}
	}

	produceWait(ctx, execCod, qtdEvs, func() int {
		return int(action.MonitorRepoCreate.Get(action.AcExecutions))
	})

	action.ForceFlush(ctx)
	showFinal(execCod, threads, qtdEvs, batchSize, time.Since(start), action.MonitorRepoCreate)
	action.Clean()

	return nil
}

func query(db *gorm.DB, execCod string, threads int, qtdQueries int, qtdEvents int, selectCod string, start string, end string) error {
	fmt.Printf("*****     QUERY Queries[%d] Events[%d] Threads[%d]\n", qtdQueries, qtdEvents, threads)
	startTm := time.Now()
	ctx := context.Background()
	for i := 0; i < qtdQueries; i++ {
		switch selectCod {
		case "RTCount":
			action.Dispatch(action.QueryRTCount(simul.GenEventNames(qtdEvents), start, end))
		case "RTSum":
			action.Dispatch(action.QueryRTSum(simul.GenEventNames(qtdEvents), start, end))
		case "RTSumValue":
			action.Dispatch(action.QueryRTSumValue(simul.GenEventNames(qtdEvents), start, end))
		case "MRTCount":
			action.Dispatch(action.QueryMRTCount(simul.GenEventNames(qtdEvents), start, end))
		case "MRTSum":
			action.Dispatch(action.QueryMRTSum(simul.GenEventNames(qtdEvents), start, end))

		default:
			action.Dispatch(action.QueryRTCount(simul.GenEventNames(qtdEvents), start, end))

		}
	}

	produceWait(ctx, execCod, qtdQueries, func() int {
		return int(action.MonitorRepoQuery.Get(action.AcExecutions))
	})
	showFinal(execCod, threads, qtdQueries, 0, time.Since(startTm), action.MonitorRepoQuery)
	action.Clean()

	return nil
}

func produceCreate(ctx context.Context, sparse bool, idxProducer int, totalPerProducer int, maxCorrelations int) {
	fnNewContext := func(bctx context.Context) context.Context {
		return simul.UserContext(
			simul.CorrelateContext(bctx, uuid.NewString()),
			simul.GenUserID(),
		)

	}

	pctx := fnNewContext(ctx)
	for ii := 0; ii < totalPerProducer; ii++ {
		action.Dispatch(action.Create(simul.NewEvent(pctx)))
		if ii%maxCorrelations == 0 {
			pctx = fnNewContext(ctx)
		}
	}
}

func produceSparse(ctx context.Context, sparse bool, idxProducer int, totalPerProducer int, maxCorrelations int) {
	fnNewContext := func(bctx context.Context) context.Context {
		return simul.UserContext(
			simul.CorrelateContext(bctx, uuid.NewString()),
			simul.GenUserID(),
		)

	}

	MaxDimensions := 1000
	SparseRatio := 0.0
	event := simul.NewEvent(fnNewContext(ctx))
	currentCorrelation := 0
	for ii := 0; ii < totalPerProducer; ii++ {
		if ii%maxCorrelations == 0 {
			event = simul.NewEvent(fnNewContext(ctx))
			currentCorrelation = 0
		}
		action.Dispatch(action.Create(event.Clone(currentCorrelation), simul.GenVector(sparse, MaxDimensions, SparseRatio)...))
		currentCorrelation++
	}
}

func show(cod string) {
	fmt.Printf("SIMUL_%s[Dispatches: %d Execs: %d Avg: %d]  \n", cod, action.MonitorActionDispatch.Get(action.AcDispatches), action.MonitorActionCreate.Get(action.AcExecutions), action.MonitorRepoCreate.Get(action.AcAvg))
}

func showFinal(cod string, threads, qtd, batch int, duration time.Duration, monitor *action.Accumulator) {
	fmt.Printf("\n\n*** SIMUL_%s[Trhreads: %d Qtd: %d Batch: %d] DURATION: %f \n", cod, threads, qtd, batch, duration.Minutes())
	fmt.Printf("*** SIMUL_%s[%s] \n\n", cod, monitor)
}

func produceWait(ctx context.Context, execCod string, qtdEvs int, getter func() int) {
	stop := false
	stalled := 0
	stalledCount := 0
	for !stop && getter() < qtdEvs {
		time.Sleep(15 * time.Second)
		if stalledCount == getter() {
			stalled++
		} else {
			stalled = 0
			stalledCount = getter()
		}
		if stalled == 10 {
			fmt.Printf("\nSTALLED\n")
			stop = true
		}
		action.ForceFlush(ctx)
		show(execCod)
	}

}

func config() {
	rand.Seed(time.Now().UnixMilli())
	time.Local = time.UTC
	os.Setenv("TZ", "UTC")
}
