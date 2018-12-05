package main

import (
	"bufio"
	"io"
	"math"
	"os"
	"strings"
)

type counter map[string]int
type weighter map[string]float64

func main() {
	contexts := loadContexts()

	numCtx := len(contexts)
	avgLen := buildAvgLen(contexts)
	ctxFreq := buildContextFrequencies(contexts)
	idf := buildIDF(numCtx, ctxFreq)
}

func loadContexts() map[string]counter {
	reader := bufio.NewReader(os.Stdin)
	contexts := make(map[string]counter)
	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		fields := strings.Fields(s)

		target := fields[0]
		ctx, ok := contexts[target]
		if !ok {
			ctx = make(counter)
			contexts[target] = ctx
		}

		for term := range fields[1:] {
			ctx[term]++
		}
	}
	return contexts
}

func buildContextFrequencies(contexts map[string]counter) counter {
	ctxFreq := make(counter)
	for _, ctx := range contexts {
		for term := range ctx {
			counter[term]++
		}
	}
	return ctxFreq
}

func buildIDF(numCtx int, ctxFreq counter) {
	idf := make(weighter)
	for term, cnt := range ctxFreq {
		idf[term] = math.Log(float64(numCtx+1) / float64(cnt))
	}
	return idf
}

func buildAvgLen(contexts map[string]counter) float64 {
	total := 0
	for _, ctx := range contexts {
		for _, cnt := range ctx {
			total += cnt
		}
	}
	numCtx := len(contexts)
	return float64(total) / float64(numCtx)
}
