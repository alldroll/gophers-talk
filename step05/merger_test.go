package main

import (
	"testing"
)

func BenchmarkMergeCandidate(b *testing.B) {
	m := NewMergeCandidate(1, 0)

	for i := 0; i < b.N; i++ {
		m.increment()
		_ = m.Position()
		_ = m.Overlap()
	}
}

type structMergeCandidate struct {
	position, overlap uint32
}

func (s structMergeCandidate) Position() DocumentID {
	return s.position
}

func (s structMergeCandidate) Overlap() int {
	return int(s.overlap)
}

func (s *structMergeCandidate) increment() {
	s.overlap++
}

func BenchmarkMergeStruct(b *testing.B) {
	m := structMergeCandidate{
		position: 1,
		overlap:  0,
	}

	for i := 0; i < b.N; i++ {
		m.increment()
		_ = m.Position()
		_ = m.Overlap()
	}
}
