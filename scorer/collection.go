package scorer

import (
	"bufio"
	"io"
	"math"
	"strings"
	"sync"
)

// Collection contains aggregated contexts for terms.
type Collection struct {
	aggContexts    *Table
	numContexts    int
	ctxFrequencies *Vector
	ctxAvgLength   float64
	idfValues      *Vector
	bm25Vectors    *Table
}

// ParadigmaticFunc is the signature of the callback for Paradigmatic function
// which is called with every pair of terms and
// their paradigmatic score.
type ParadigmaticFunc func(a, b string, score float64)

// SyntagmaticFunc is the signature of the callback for Syntagmatic function
// which is called with  every co-occurred pair of terms and
// their syntagmatic score.
type SyntagmaticFunc func(a, b string, score float64)

// NewCollection creates a collection by reading lines from reader.
// Each line is considered as a context which consists of
// terms separated by spaces. k, b are parameters for BM25 algorithm.
func NewCollection(reader *bufio.Reader, k, b float64) *Collection {
	col := &Collection{}
	col.load(reader)
	col.initIDFValues()
	col.initBM25Vectors(k, b)
	return col
}

func (col *Collection) load(reader *bufio.Reader) {
	// Read and count terms
	col.aggContexts = NewTable()
	col.ctxFrequencies = NewVector()

	ctxTotalLength := 0.0
	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		terms := strings.Fields(s)
		col.aggContexts.Update(terms)
		col.numContexts++

		for _, term := range terms {
			col.ctxFrequencies.Increment(term)
		}

		ctxTotalLength += float64(len(terms))
	}
	col.ctxAvgLength = ctxTotalLength / float64(col.numContexts)
}

func (col *Collection) initIDFValues() {
	M := float64(col.numContexts)
	calcIDF := func(k float64) float64 {
		// NOTE: Use Probabilistic IDF to make frequent terms less important.
		return math.Log((M - k + 1) / k)
	}

	col.idfValues = col.ctxFrequencies.Map(calcIDF)
}

func (col *Collection) initBM25Vectors(k, b float64) {
	avgLen := col.ctxAvgLength

	calcBM25 := func(aggContext *Vector) *Vector {
		ctxLen := aggContext.Total()
		lenNormalizer := k * (1 - b + (b * ctxLen / avgLen))

		bm25Values := aggContext.Map(func(count float64) float64 {
			return ((k + 1) * count) / (count + lenNormalizer)
		})
		// Normalize
		bm25Vector := bm25Values.Map(func(bm25 float64) float64 {
			return bm25 / bm25Values.Total()
		})

		return bm25Vector
	}

	col.bm25Vectors = col.aggContexts.Map(calcBM25)
}

// Paradigmatic calls f on every pair of terms and
// their paradigmatic relationship score based on BM25 algorithm.
// Note that (a, b) and (b, a) are different pairs,
// even though the score of (a, b) and (b, a) are the same.
func (col *Collection) Paradigmatic(numWorkers int, f ParadigmaticFunc) {
	worker := func(w work) {
		termA := w.term
		bm25VecA := w.vector
		for termB, bm25VecB := range col.bm25Vectors.self {
			if termA == termB {
				continue
			}
			score := similarity(bm25VecA, bm25VecB, col.idfValues)
			f(termA, termB, score)
		}
	}

	col.workParallel(numWorkers, worker)
}

func similarity(a, b *Vector, idfValues *Vector) float64 {
	shorter, other := a, b
	if b.Len() < a.Len() {
		shorter, other = b, a
	}

	score := 0.0
	for term, w1 := range shorter.self {
		if w2, ok := other.self[term]; ok {
			score += idfValues.self[term] * w1 * w2
		}
	}

	return score
}

// Syntagmatic call f on every co-occurred pair and
// their syntagmatic relationship score based on BM25 algorithm.
func (col *Collection) Syntagmatic(numWorkers int, f SyntagmaticFunc) {
	worker := func(w work) {
		termA := w.term
		bm25VecA := w.vector
		for termB, bm25 := range bm25VecA.self {
			if termA == termB {
				continue
			}
			score := idfWeighted(termB, bm25, col.idfValues)
			f(termA, termB, score)
		}
	}

	col.workParallel(numWorkers, worker)
}

func idfWeighted(term string, bm25 float64, idfValues *Vector) float64 {
	return bm25 * idfValues.self[term]
}

// NumTerms returns the number of occurred terms.
func (col *Collection) NumTerms() int {
	return len(col.aggContexts.self)
}

type work struct {
	term   string
	vector *Vector
}

type workerFunc func(w work)

func (col *Collection) workParallel(numWorkers int, worker workerFunc) {
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	works := make(chan work, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for w := range works {
				worker(w)
			}
		}()
	}

	for termA, bm25VecA := range col.bm25Vectors.self {
		works <- work{termA, bm25VecA}
	}
	close(works)
	wg.Wait()
}
