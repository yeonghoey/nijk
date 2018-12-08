package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yeonghoey/nijk/scorer"
)

const (
	k = 1.2
	b = 0.75
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	collection := scorer.NewCollection(reader, k, b)

	collection.Paradigmatic(func(a, b string, score float64) {
		fmt.Printf("%s %s %.2f\n", a, b, score)
	})

	collection.Syntagmatic(func(a, b string, score float64) {
		fmt.Printf("%s %s %.2f\n", a, b, score)
	})
}
