package main

import "sort"

// DocumentID is a unique identifier of a document
type DocumentID = uint32

// Term is a search item of a document
type Term = string

// PostingList represents the list of documents that belongs to a search term
type PostingList = []DocumentID

// InvertedIndex is a datastructure, that maps back from terms to the parts of a document where they occur
type InvertedIndex = map[string]PostingList

// LengthFilter maps back from terms length (nGram count) to the inverted indexes
type LengthFilter = []InvertedIndex

// NGramIndex represents a fuzzy string search structure
type NGramIndex struct {
	nGram      int
	index      LengthFilter
	dictionary []string
}

// BuildIndex builds the inverted index structure for the given dictionary
func BuildIndex(nGram int, dictionary []string) *NGramIndex {
	index := LengthFilter{}

	for id, word := range dictionary {
		nGrams := SplitIntoNGrams(nGram, normalize(word))
		n := len(nGrams)

		if n >= len(index) {
			tmp := make(LengthFilter, n+1)
			copy(tmp, index)
			index = tmp
		}

		if index[n] == nil {
			index[n] = make(InvertedIndex)
		}

		for _, term := range nGrams {
			if _, ok := index[n][term]; !ok {
				index[n][term] = PostingList{}
			}

			index[n][term] = append(index[n][term], uint32(id))
		}
	}

	return &NGramIndex{
		nGram:      nGram,
		index:      index,
		dictionary: dictionary,
	}
}

// Search performs approximate strign search for the given query
func (n *NGramIndex) Search(query string) SearchResult {
	terms := SplitIntoNGrams(n.nGram, normalize(query))
	result := SearchResult{}
	sizeA := len(terms)
	rid := []PostingList{}

	for sizeB, index := range n.index {
		if index == nil {
			continue
		}

		for _, term := range terms {
			posting := index[term]

			if posting == nil {
				continue
			}

			rid = append(rid, posting)
		}

		for _, candidate := range Merge(rid) {
			docID := candidate.Position()

			result = append(result, SearchCandidate{
				Candidate: n.dictionary[int(docID)],
				Score:     JaccardDistance(candidate.Overlap(), sizeA, sizeB),
			})
		}

		rid = rid[:0]
	}

	sort.Stable(result)

	return result
}

// SplitIntoNGrams splits the given query on nGrams
func SplitIntoNGrams(nGram int, query string) []Term {
	runes := []rune(query)

	if len(runes) < nGram {
		return []Term{}
	}

	result := make([]Term, 0, len(runes)-nGram+1)

	for i := 0; i < len(runes)-nGram+1; i++ {
		result = appendUnique(result, string(runes[i:i+nGram]))
	}

	return result
}

// appendUnique appends an item only to slice if there is not such item
func appendUnique(slice []Term, item Term) []Term {
	for _, c := range slice {
		if c == item {
			return slice
		}
	}

	return append(slice, item)
}
