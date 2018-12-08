package scorer

// Table represents vectors associated with terms.
type Table struct {
	self map[string]*Vector
}

// TableMapFunc is the signature of Map function used to
// create a new mapped Table based on the original one.
type TableMapFunc func(x *Vector) (y *Vector)

// TableEachFunc is the signature of Each function used to
// run an arbitrary function on each term, vector pair.
type TableEachFunc func(term string, v *Vector)

// NewTable returns a new empty Table
func NewTable() *Table {
	return &Table{map[string]*Vector{}}
}

// Get returns the associated vector for term.
func (t *Table) Get(term string) *Vector {
	vector, ok := t.self[term]
	if !ok {
		vector = NewVector()
		t.self[term] = vector
	}
	return vector
}

// Update updates each vector of terms by incrementing values of the other terms.
func (t *Table) Update(terms []string) {
	for _, term := range terms {
		vector := t.Get(term)
		for _, other := range terms {
			// Should not counter the term itself
			if term == other {
				continue
			}
			vector.Increment(other)
		}
	}
}

// Map returns a new Table with vectors mapaed by f.
func (t *Table) Map(f TableMapFunc) *Table {
	tm := NewTable()
	for term, x := range t.self {
		tm.self[term] = f(x)
	}
	return tm
}

// Each calls f on each term, vector pair.
func (t *Table) Each(f TableEachFunc) {
	for term, vector := range t.self {
		f(term, vector)
	}
}
