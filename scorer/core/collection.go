package core

import (
	"bufio"
	"io"
	"math"
	"strings"
)

// Collection contains aggregated contexts for terms.
type Collection struct {
	terms          []string
	aggContexts    *Table
	numContexts    int
	ctxFrequencies *Vector
	ctxAvgLength   float64
	idfValues      *Vector
	bm25Vectors    *Table
}

// ParadigmaticFunc is the signature of Paradigmatic function used as a callback
// for every pair of terms and their paradigmatic score.
type ParadigmaticFunc func(a, b string, score float64)

// NewCollection creates a collection by reading lines from reader.
// Each line is considered as a context which consists of terms separated by spaces.
// k, b are parameters for BM25 algorithm.
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

	termsExisting := map[string]bool{}
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

			if !termsExisting[term] {
				col.terms = append(col.terms, term)
				termsExisting[term] = true
			}
		}

		ctxTotalLength += float64(len(terms))
	}
	col.ctxAvgLength = ctxTotalLength / float64(col.numContexts)
}

func (col *Collection) initIDFValues() {
	M := float64(col.numContexts)
	calcIDF := func(k float64) float64 {
		return math.Log((M + 1) / k)
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
		normalized := bm25Values.Map(func(bm25 float64) float64 {
			return bm25 / bm25Values.Total()
		})

		return normalized
	}

	col.bm25Vectors = col.aggContexts.Map(calcBM25)
}

// Paradigmatic calls f on every distinct pair of terms and their paradigmatic relationship score
// based on BM25 algorithm.
func (col *Collection) Paradigmatic(f ParadigmaticFunc) {
	for ai := 0; ai < len(col.terms); ai++ {
		termA := col.terms[ai]
		bm25A := col.bm25Vectors.Get(termA)
		for bi := ai + 1; bi < len(col.terms); bi++ {
			termB := col.terms[bi]
			bm25B := col.bm25Vectors.Get(termB)
			score := similarity(bm25A, bm25B, col.idfValues)
			f(termA, termB, score)
		}
	}
}

func similarity(bm25A, bm25B *Vector, idfValues *Vector) float64 {
	shorter, other := bm25A, bm25B
	if bm25B.Len() < bm25A.Len() {
		shorter, other = bm25B, bm25A
	}

	score := 0.0
	for term, w1 := range shorter.self {
		if w2, ok := other.self[term]; ok {
			score += idfValues.self[term] * w1 * w2
		}
	}

	return score
}
