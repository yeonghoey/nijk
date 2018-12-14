package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"

	"github.com/yeonghoey/nijk/scorer"
)

// TODO: Parameterize these constants
const (
	numWorkers = 64

	k = 1.2
	b = 0.75

	paradigmaticThreshold = 0.5
	syntagmaticThreshold  = 0.5
)

var preset string = os.Args[1]

func main() {
	reader := bufio.NewReader(os.Stdin)
	collection := scorer.NewCollection(reader, k, b)

	log.Printf("Run Paradigmatic")
	collection.Paradigmatic(numWorkers, newHandler("paradigmatic", paradigmaticThreshold))

	log.Printf("Run Syntagmatic")
	collection.Syntagmatic(numWorkers, newHandler("syntagmatic", syntagmaticThreshold))
}

func newHandler(relation string, threshold float64) func(a, b string, score float64) {
	var mutex = &sync.Mutex{}
	var processed int32
	return func(a, b string, score float64) {
		if score < threshold {
			return
		}
		mutex.Lock()
		fmt.Println(insertQuery(preset, relation, a, b, score))
		mutex.Unlock()

		incremented := atomic.AddInt32(&processed, 1)
		if incremented%1000 == 0 {
			log.Printf("%d processed", incremented)
		}
	}

}

func insertQuery(preset, relation string, this, that string, score float64) string {
	table := fmt.Sprintf("`%s_%s`", preset, relation)
	return fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", %.5f);",
		table, this, that, score)
}
