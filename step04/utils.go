package main

import "strings"

// normalize performs normalization for the given query
func normalize(query string) string {
	q := strings.ToLower(query)
	q = strings.TrimSpace(q)

	return "^" + q + "$"
}
