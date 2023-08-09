package simul

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	grand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

var dist = distuv.Normal{
	Mu:    500, // Mean of the normal distribution
	Sigma: 50,  // Standard deviation of the normal distribution
	Src:   grand.NewSource(uint64(time.Now().UnixNano())),
}

func genEventName() string {
	index := uint64(dist.Rand())
	if index > 1000 {
		index = 1000
	}
	return fmt.Sprintf("ev_bus_%d", index)
}

func genPayload(eventName string) string {
	return fmt.Sprintf("{\"value\": %f}", math.Trunc(math.Abs(rand.NormFloat64()*100)))
}

func genLabels(eventName string) string {
	fGenLabel := func() []string {
		lbs := []string{}
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		qtd := rnd.Intn(4)
		for i := 0; i < qtd; i++ {
			rnd.Seed(time.Now().UnixNano())
			lbs = append(lbs, fmt.Sprintf("bus%d", rnd.Intn(10)+1))
		}
		return lbs
	}
	labels := []string{"operational"}
	labels = append(labels, fGenLabel()...)

	return strings.Join(labels, ",")
}
