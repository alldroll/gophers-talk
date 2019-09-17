package main

// MergeCandidate is a pair <DocumentID, overlap count>
type MergeCandidate uint64

// NewMergeCandidate creates a new instance of MergeCandidate for the given document id and overlap count
func NewMergeCandidate(id DocumentID, overlap uint32) MergeCandidate {
	return MergeCandidate(pack(id, overlap))
}

// Position returns the position of the merge candidate
func (m MergeCandidate) Position() DocumentID {
	docID, _ := unpack(uint64(m))

	return docID
}

// Overlap returns the overlap count of the merge candidate
func (m MergeCandidate) Overlap() int {
	_, overlap := unpack(uint64(m))

	// This will be safe, because we will never reach overlap (we don't have so many posting list at once)
	return int(overlap)
}

// increment increments the underlying overlap count
func (m *MergeCandidate) increment() {
	*m = NewMergeCandidate(m.Position(), uint32(m.Overlap()+1))
}

// Merge merges the given list of postings and returns array of merge candidates
// The returned list holds merge candidates with preserving the order
func Merge(rid []PostingList) []MergeCandidate {
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

// pack packes the given 2 uint32 numbers into uint64
func pack(a, b uint32) uint64 {
	return uint64(a)<<32 | uint64(b)
}

// unpack explodes the given uint64 into 2 uint32
func unpack(m uint64) (uint32, uint32) {
	return uint32(m >> 32), uint32(m)
}
