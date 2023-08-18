package main

import (
	"context"
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
	SIZE := 100000000
	start := time.Now()
	values := make([]string, SIZE)
	count := 0
	for i := 0; i < SIZE; i++ {
		values[i] = fmt.Sprintf("%d", i)
		count++
	}
	sec := time.Since(start).Seconds()
	fmt.Printf("Tempo: %f\n", sec)
	fmt.Printf("Size: %d\n", len(values))
}

func main() {
	var codFlag = flag.String("cod", "0", "simul code")
	var createFlag = flag.Bool("create", false, "Create events")
	var threadsFlag = flag.Int("threads", 4, "Paralel instances")
	var qtdFlag = flag.Int("qtd", 10, "Qtd of events")
	var batchFlag = flag.Int("batch", 10, "Qtd batch")
	var deamonFlag = flag.Bool("deamon", false, "Continuous run")
	var queryFlag = flag.String("query", "NRT", "Query Code")

	flag.Parse()

	config()

	db, err := gorm.Open(mysql.Open("root:singlestore@tcp(10.164.47.110:3306)/events?parseTime=true&loc=UTC"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
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
		BatchSize:   *batchFlag,
		DB:          db,
	})

	if *createFlag {
		if *deamonFlag {
			for i := 0; ; i++ {
				fmt.Println("\n\n***")
				fmt.Printf("*** SIMUL_%s EXECUTION %d \n", *codFlag, i)
				fmt.Println("***")
				create(db, *codFlag, *threadsFlag, *qtdFlag, *batchFlag)
			}
		} else {
			create(db, *codFlag, *threadsFlag, *qtdFlag, *batchFlag)
		}
	} else {
		query(db, *codFlag, *threadsFlag, *qtdFlag, *batchFlag, *queryFlag)
	}
}

func create(db *gorm.DB, execCod string, threads int, qtdEvs int, batchSize int) error {
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
		go produceCreate(ctx, i, count)
	}

	produceWait(ctx, execCod, qtdEvs, func() int {
		return int(action.MonitorRepoCreate.Get(action.AcExecutions))
	})

	action.ForceFlush(ctx)
	showFinal(execCod, threads, qtdEvs, batchSize, time.Since(start), action.MonitorRepoCreate)
	action.Clean()

	return nil
}

func query(db *gorm.DB, execCod string, threads int, qtdQueries int, qtdEvents int, query string) error {
	start := time.Now()
	ctx := context.Background()
	for i := 0; i < qtdQueries; i++ {
		switch query {
		case "RTCount":
			action.Dispatch(action.QueryRTCount(simul.GenEventNames(qtdEvents), "2023-08-09 12:00:00", "2023-08-13 19:00:00"))
		case "RTSum":
			action.Dispatch(action.QueryRTSum(simul.GenEventNames(qtdEvents), "2023-08-09 12:00:00", "2023-08-13 19:00:00"))
		default:
			action.Dispatch(action.QueryRTCount(simul.GenEventNames(qtdEvents), "2023-08-09 12:00:00", "2023-08-13 19:00:00"))

		}
	}

	produceWait(ctx, execCod, qtdQueries, func() int {
		return int(action.MonitorRepoQuery.Get(action.AcExecutions))
	})
	showFinal(execCod, threads, qtdQueries, 0, time.Since(start), action.MonitorRepoQuery)
	return nil
}

func produceCreate(ctx context.Context, idx int, total int) {
	MaxCorrelations := 100
	fnNewContext := func(bctx context.Context) context.Context {
		return simul.UserContext(
			simul.CorrelateContext(bctx, uuid.NewString()),
			simul.GenUserID(),
		)

	}

	pctx := fnNewContext(ctx)

	for ii := 0; ii < total; ii++ {
		action.Dispatch(action.Create(simul.NewEvent(pctx)))
		if ii%MaxCorrelations == 0 {
			pctx = fnNewContext(ctx)
		}
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
		if stalledCount == int(action.MonitorRepoCreate.Get(action.AcExecutions)) {
			stalled++
		} else {
			stalled = 0
			stalledCount = int(action.MonitorRepoCreate.Get(action.AcExecutions))
		}
		if stalled == 4 {
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
