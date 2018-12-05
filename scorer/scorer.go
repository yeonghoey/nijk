package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

type counter map[string]int
type weighter map[string]float64

type context struct {
	termFreq counter
	length   int
}

type record struct {
	other string
	score float64
}

type byScoreDesc []record

func (a byScoreDesc) Len() int           { return len(a) }
func (a byScoreDesc) Less(i, j int) bool { return a[i].score > a[j].score }
func (a byScoreDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

const (
	k = 1.2
	b = 0.75
)

func main() {
	contexts := loadContexts()

	ctxFreq := buildCtxFreq(contexts)
	idf := buildIDF(contexts, ctxFreq)

	avgLen := calcAvgLen(contexts)
	bm25 := buildBM25(contexts, avgLen)

	records := buildRecords(bm25, idf)
	for target, rs := range records {
		for _, r := range rs {
			fmt.Printf("%s %s %f\n", target, r.other, r.score)
		}
	}
}

func loadContexts() map[string]*context {
	reader := bufio.NewReader(os.Stdin)
	contexts := map[string]*context{}
	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		fields := strings.Fields(s)

		target := fields[0]
		ctx, ok := contexts[target]
		if !ok {
			ctx = &context{counter{}, 0}
			contexts[target] = ctx
		}

		for _, term := range fields[1:] {
			ctx.termFreq[term]++
			ctx.length++
		}
	}
	return contexts
}
func calcAvgLen(contexts map[string]*context) float64 {
	total := 0
	for _, ctx := range contexts {
		total += ctx.length
	}
	numCtx := len(contexts)
	return float64(total) / float64(numCtx)
}

func buildBM25(contexts map[string]*context, avgLen float64) map[string]weighter {
	calc := func(target string, ctx *context) weighter {
		lenNormalizer := k * (1 - b + (b*float64(ctx.length))/avgLen)
		total := 0.0
		terms := weighter{}
		for term, cnt := range ctx.termFreq {
			score := ((k * 1) * float64(cnt)) / (float64(cnt) + lenNormalizer)
			terms[term] = score
			total += score
		}

		for term, score := range terms {
			terms[term] = score / total
		}
		return terms
	}

	bm25 := map[string]weighter{}
	for target, ctx := range contexts {
		bm25[target] = calc(target, ctx)
	}
	return bm25
}

func buildCtxFreq(contexts map[string]*context) counter {
	ctxFreq := counter{}
	for _, ctx := range contexts {
		for term := range ctx.termFreq {
			ctxFreq[term]++
		}
	}
	return ctxFreq
}

func buildIDF(contexts map[string]*context, ctxFreq counter) weighter {
	numCtx := len(contexts)
	idf := weighter{}
	for term, cnt := range ctxFreq {
		idf[term] = math.Log(float64(numCtx+1) / float64(cnt))
	}
	return idf
}

func buildRecords(bm25 map[string]weighter, idf weighter) map[string][]record {
	records := map[string][]record{}
	for target, bm1 := range bm25 {
		for other, bm2 := range bm25 {
			if target == other {
				continue
			}
			score := similarity(bm1, bm2, idf)
			if score < 0.1 {
				continue
			}
			records[target] = append(records[target], record{other, score})
		}
		sort.Sort(byScoreDesc(records[target]))
	}
	return records
}

func similarity(a, b weighter, idf weighter) float64 {
	shorter, other := a, b
	if len(b) < len(a) {
		shorter, other = b, a
	}

	score := 0.0
	for term, w1 := range shorter {
		if w2, ok := other[term]; ok {
			score += idf[term] * w1 * w2
		}
	}
	return score
}
