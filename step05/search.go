package main

// SearchCandidate represents the result item
type SearchCandidate struct {
	candidate MergeCandidate
	Value     string
	Score     float64
}

// less tells if the given candidate has weaker characteristics than the provided
func (s SearchCandidate) less(candidate MergeCandidate, score float64) bool {
	if s.candidate.Overlap() == candidate.Overlap() {
		return s.Score > score
	}

	return s.candidate < candidate
}

// SearchResult is a result of fuzzy search request
type SearchResult []SearchCandidate

// Len is the number of elements in the collection.
func (s SearchResult) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SearchResult) Less(i, j int) bool {
	return s[i].less(s[j].candidate, s[j].Score)
}

// Swap swaps the elements with indexes i and j.
func (s SearchResult) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Push pushes the given x to the collection list
func (s *SearchResult) Push(x interface{}) {
	*s = append(*s, x.(SearchCandidate))
}

// Pop pops the item from the heap
func (s *SearchResult) Pop() interface{} {
	old := *s
	n := len(old)
	x, old := old[n-1], old[:n-1]
	*s = old

	return x
}
