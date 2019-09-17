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

// NGramIndex represents a fuzzy string search structure
type NGramIndex struct {
	nGram      int
	index      InvertedIndex
	dictionary []string
}

// BuildIndex builds the inverted index structure for the given dictionary
func BuildIndex(nGram int, dictionary []string) *NGramIndex {
	index := make(InvertedIndex)

	for id, word := range dictionary {
		for _, term := range SplitIntoNGrams(nGram, word) {
			if _, ok := index[term]; !ok {
				index[term] = PostingList{}
			}

			index[term] = append(index[term], uint32(id))
		}
	}

	return &NGramIndex{
		nGram:      nGram,
		index:      index,
		dictionary: dictionary,
	}
}

// Search performs approximate strign search for the given query
func (n *NGramIndex) Search(query string) []string {
	terms := SplitIntoNGrams(n.nGram, query)
	result := []string{}
	rid := []PostingList{}

	for _, term := range terms {
		posting := n.index[term]

		if posting == nil {
			continue
		}

		rid = append(rid, posting)
	}

	candidates := Merge(rid)

	// sort by overlap count
	// use SliceStable to preserve the order for the candidates with the same overlap count
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Overlap() > candidates[j].Overlap()
	})

	for _, candidate := range candidates {
		docID := candidate.Position()
		result = append(result, n.dictionary[int(docID)])
	}

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
