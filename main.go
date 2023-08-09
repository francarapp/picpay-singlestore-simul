package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/francarapp/picpay-singlestore-simul/pkg/action"
	"github.com/francarapp/picpay-singlestore-simul/pkg/domain"
	"github.com/francarapp/picpay-singlestore-simul/pkg/simul"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	var createFlag = flag.Bool("create", true, "Create events")
	var threadsFlag = flag.Int("threads", 2, "Paralel instances")
	var qtdFlag = flag.Int("qtd", 100, "Qtd of events")
	var batchFlag = flag.Int("batch", 10, "Qtd batch")

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
		create(db, *threadsFlag, *qtdFlag, *batchFlag)
	} else {
		query(db, *threadsFlag)
	}
}

func create(db *gorm.DB, threads int, qtd int, batch int) error {
	ctx := context.Background()
	fnNewContext := func(bctx context.Context) context.Context {
		return context.WithValue(bctx, simul.CorrelationKey, simul.CorrelationID(uuid.New().String()))
	}
	producers := threads - 1
	if producers == 0 {
		producers = 1
	}
	for i := 0; i < producers; i++ {
		go func() {
			MaxCorrelations := 100
			pctx := fnNewContext(ctx)
			count := int(qtd / producers)
			if qtd%producers != 0 && i == 0 {
				count += qtd % producers
			}
			for ii := 0; ii < count; ii++ {
				action.Dispatch(action.Create(simul.NewEvent(pctx)))
				if ii%MaxCorrelations == 0 {
					pctx = fnNewContext(ctx)
				}
			}
		}()
	}

	stop := false
	stalled := 0
	stalledCount := 0
	for !stop && action.Monitor.Creations < qtd {
		time.Sleep(15 * time.Second)
		if stalledCount == action.Monitor.Creations {
			stalled++
		} else {
			stalled = 0
			stalledCount = action.Monitor.Creations
		}
		if stalled == 4 {
			stop = true
		}
		action.Flush(ctx)
		fmt.Printf("Dispatches: %d Creates: %d %v", action.MonitorDispatch.Get(), action.MonitorCreate.Get(), action.Monitor)
	}

	action.Flush(ctx)
	fmt.Printf("Dispatches: %d Creates: %d %v", action.MonitorDispatch.Get(), action.MonitorCreate.Get(), action.Monitor)
	return nil
}

func query(db *gorm.DB, instances int) error {
	var tx *gorm.DB
	var events []domain.Event
	db.Where("event_name = ?", "bus_ev_1").Find(&events)
	fmt.Printf("Events: %s", events)

	var results []map[string]interface{}
	tx = db.Model(&domain.Event{}).Select("payload::$key_a").Where("event_name = ?", "bus_ev_1").Find(&results)
	if tx.Error != nil {
		fmt.Printf("Failed: %s", tx.Error)
	} else {
		fmt.Printf("Payload: %s", results)
	}

	results = make([]map[string]interface{}, 0)
	tx = db.Raw("select payload::$key_a as key_a from event where event_name = ?", "bus_ev_1").Find(&results)
	if tx.Error != nil {
		fmt.Printf("Failed: %s", tx.Error)
	} else {
		fmt.Printf("Payload: %s", results)
	}
	return nil
}

func config() {
	rand.Seed(time.Now().UnixMilli())
	time.Local = time.UTC
	os.Setenv("TZ", "UTC")
}
