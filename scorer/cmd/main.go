// Package main runs scorer on os.Stdin and prints SQLs to Os.Stdout for generating a dump file of the analysis result.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/yeonghoey/nijk/scorer"
)

// TODO: Parameterize these constants
const (
	numWorkers = 64
	topN       = 50

	k = 1.2
	b = 0.75

	paradigmaticThreshold = 0.5
	syntagmaticThreshold  = 0
)

type entry struct {
	that  string
	score float64
}

type scoreDesc []entry

func (a scoreDesc) Len() int           { return len(a) }
func (a scoreDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a scoreDesc) Less(i, j int) bool { return a[i].score > a[j].score }

var preset string = os.Args[1]

func main() {
	reader := bufio.NewReader(os.Stdin)
	collection := scorer.NewCollection(reader, k, b)

	log.Printf("Run Paradigmatic")

	parEntries := map[string][]entry{}
	collection.Paradigmatic(numWorkers, newHandler(parEntries, paradigmaticThreshold))
	outputEntries("paradigmatic", parEntries)

	log.Printf("Run Syntagmatic")
	synEntries := map[string][]entry{}
	collection.Syntagmatic(numWorkers, newHandler(synEntries, syntagmaticThreshold))
	outputEntries("syntagmatic", synEntries)
}

func newHandler(entries map[string][]entry, threshold float64) func(a, b string, score float64) {
	var mutex = &sync.Mutex{}
	var processed int32
	return func(a, b string, score float64) {
		if score < threshold {
			return
		}
		mutex.Lock()
		entries[a] = append(entries[a], entry{b, score})
		mutex.Unlock()

		incremented := atomic.AddInt32(&processed, 1)
		if incremented%1000 == 0 {
			log.Printf("%d processed", incremented)
		}
	}

}

func outputEntries(relation string, entries map[string][]entry) {
	for this, es := range entries {
		sort.Sort(scoreDesc(es))
		for _, e := range es[:min(len(es), topN)] {
			fmt.Println(insertQuery(relation, this, e.that, e.score))
		}
	}
}

func insertQuery(relation string, this, that string, score float64) string {
	table := fmt.Sprintf("`%s_%s`", preset, relation)
	return fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", %.5f);",
		table, this, that, score)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
