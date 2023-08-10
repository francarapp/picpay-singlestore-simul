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

var distJsonKeys = distuv.Normal{
	Mu:    20, // Mean of the normal distribution
	Sigma: 5,  // Standard deviation of the normal distribution
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
	q := int(distJsonKeys.Rand())
	if q == 0 {
		q = 1
	} else if q > 40 {
		q = 40
	}
	json := []string{
		fmt.Sprintf("\"value\": %f", math.Trunc(math.Abs(rand.NormFloat64()*100))),
	}
	for i := 1; i <= q; i++ {
		json = append(json,
			fmt.Sprintf("\"key_%d\":\"%s\"", i,
				fmt.Sprintf("xpto_%d", rand.Intn(100)+1),
			),
		)
	}
	return fmt.Sprintf("{%s}", strings.Join(json, ","))
}

func genLabels(eventName string) string {
	fGenLabel := func() []string {
		lbs := []string{}
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		qtd := rnd.Intn(10)
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
