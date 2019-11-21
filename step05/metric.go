package main

// JaccardSimilarity calculates the similarity between 2 sets of nGrams
func JaccardSimilarity(inter, sizeA, sizeB int) float64 {
	return float64(inter) / float64(sizeA+sizeB-inter)
}

// JaccardDistance returns the distance between 2 sets of nGrams
func JaccardDistance(inter, sizeA, sizeB int) float64 {
	return 1 - JaccardSimilarity(inter, sizeA, sizeB)
}
