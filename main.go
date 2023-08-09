package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/francarapp/picpay-singlestore-simul/pkg/action"
	"github.com/francarapp/picpay-singlestore-simul/pkg/domain"
	"github.com/francarapp/picpay-singlestore-simul/pkg/simul"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		}},
	)
	if err != nil {
		panic("failed to connect database")
	}

	action.InitDispatching(&action.DispatchConfig{
		ChanSize:    100000,
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
	fnNewContext := func() context.Context {
		return context.WithValue(context.Background(), simul.CorrelationKey, simul.CorrelationID(uuid.New().String()))
	}
	MaxCorrelations := 100
	ctx := fnNewContext()
	for i := 0; i < qtd; i++ {
		action.Dispatch(action.Create(simul.NewEvent(ctx)))
		if i%MaxCorrelations == 0 {
			ctx = fnNewContext()
		}
	}
	for action.Monitor.Creations < qtd {
		time.Sleep(60 * time.Second)
	}
	fmt.Printf("%v", action.Monitor)
	action.Flush(ctx)
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
}
