/*
Package scorer analyzes a collection to discover paradigmatic and syntagmatic relation.

  collection := scorer.NewCollection(reader, k, b)
  collection.Paradigmatic(func (a, b string, score float64) {
    // Do whatever you want with a, b and their paradigmatic score.
  })
  collection.Syntagmatic(func (a, b string, score float64) {
    // Do whatever you want with a, b and their paradigmatic score.
  })

The reader provided to scorer.NewCollection should provide contexts by space-separated lines.
The collection aggregates the term counts for each term of contexts to form aggregated contexts.
The aggregated context of each term is basically the basic unit of the term's context to be used for determining the similarity between two terms.
Also, the aggregated context of a term is used for the term's syntagmatic analysis, by element-wise multiplication of IDF.

The analysis is based on BM25 algorithm. Since the context generally has a few terms, the implementation uses a hash map to represent a context vector,
instead of an array of the vocabulary size.
*/
package scorer
