package main

// MergeCandidate is a pair <DocumentID, overlap count>
type MergeCandidate uint64

// NewMergeCandidate creates a new instance of MergeCandidate for the given document id and overlap count
func NewMergeCandidate(id DocumentID, overlap int) MergeCandidate {
	return MergeCandidate(uint64(overlap)<<32 | uint64(id))
}

// Position returns the position of the merge candidate
func (m MergeCandidate) Position() DocumentID {
	return DocumentID(uint32(m))
}

// Overlap returns the overlap count of the merge candidate
func (m MergeCandidate) Overlap() int {
	// This will be safe, because we will never reach overlap (we don't have so many posting list at once)
	return int(m >> 32)
}

// increment increments the underlying overlap count
func (m *MergeCandidate) increment() {
	*m = *m + 1<<32
}

// Merge merges the given list of postings and returns array of merge candidates
// The returned list holds merge candidates with preserving the order
func Merge(rid []PostingList) []MergeCandidate {
	// use an approach of merging 2 sorted list described here https://www.geeksforgeeks.org/merge-two-sorted-arrays/
	// here tmp will store the result of the merge between (rid[k] and result)
	// after merge operation, we swap tmp and result and do the same operation with rid[k+1]

	result := []MergeCandidate{}
	tmp := []MergeCandidate{}

	for _, list := range rid {
		i, j := 0, 0
		endI, endJ := len(list), len(result)

		for j < endJ || i < endI {
			if j >= endJ || (i < endI && list[i] < result[j].Position()) {
				tmp = append(tmp, NewMergeCandidate(list[i], 1))
				i++
			} else if i >= endI || j < endJ && result[j].Position() < list[i] {
				tmp = append(tmp, result[j])
				j++
			} else {
				result[j].increment()
				tmp = append(tmp, result[j])
				i++
				j++
			}
		}

		tmp, result = result, tmp
		tmp = tmp[:0]
	}

	return result
}
