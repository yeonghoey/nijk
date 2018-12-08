package core

// Vector represents values, which are intended to be
// frequencies or weights, of terms.
type Vector struct {
	self  map[string]float64
	total float64
}

// VectorMapFunc is the signature of Map function used to
// create a new mapped Vector based on the original one.
type VectorMapFunc func(x float64) (y float64)

// VectorEachFunc is the signature of Each function used to
// run an arbitrary function on each term, value pair.
type VectorEachFunc func(term string, x float64)

// NewVector returns a new empty Vector.
func NewVector() *Vector {
	return &Vector{map[string]float64{}, 0.0}
}

// Len returns the number of value-existing terms.
func (v *Vector) Len() int {
	return len(v.self)
}

// Total returns the sum of values.
func (v *Vector) Total() float64 {
	return v.total
}

// Increment increments the value of the term by 1.
func (v *Vector) Increment(term string) {
	v.self[term] += 1.0
	v.total += 1.0
}

// Map returns a new Vector with values mapped by f.
func (v *Vector) Map(f VectorMapFunc) *Vector {
	vm := NewVector()
	for term, x := range v.self {
		y := f(x)
		vm.self[term] = y
		vm.total += y
	}
	return vm
}

// Each calls f on each term value pair.
func (v *Vector) Each(f VectorEachFunc) {
	for term, x := range v.self {
		f(term, x)
	}
}
