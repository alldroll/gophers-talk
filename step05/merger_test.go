package main

import (
	"sort"
	"testing"
)

func BenchmarkMergeCandidate(b *testing.B) {
	m := NewMergeCandidate(1, 0)

	for i := 0; i < b.N; i++ {
		m.increment()

		for j := 0; j < 1000; j++ {
			_ = m.Position()
			_ = m.Overlap()
		}
	}
}

func BenchmarkMergeStruct(b *testing.B) {
	m := structMergeCandidate{
		position: 1,
		overlap:  0,
	}

	for i := 0; i < b.N; i++ {
		m.increment()

		for j := 0; j < 1000; j++ {
			_ = m.Position()
			_ = m.Overlap()
		}
	}
}

func BenchmarkMergeCandidateSort(b *testing.B) {
	b.StopTimer()

	m := make([]MergeCandidate, 0, 1<<16)
	data := make([]MergeCandidate, len(m))

	for i := DocumentID(0); i < 10000; i++ {
		overlap := i ^ 0xcccc
		m = append(m, NewMergeCandidate(i, int(overlap)))
	}

	for i := 0; i < b.N; i++ {
		copy(data, m)
		b.StartTimer()

		sort.SliceStable(data, func(i, j int) bool {
			return data[i] > data[j]
		})

		b.StopTimer()
	}
}

func BenchmarkMergeStructSort(b *testing.B) {
	b.StopTimer()

	m := make([]structMergeCandidate, 0, 1<<16)
	data := make([]structMergeCandidate, len(m))

	for i := DocumentID(0); i < 10000; i++ {
		overlap := i ^ 0xcccc

		m = append(m, structMergeCandidate{
			position: i,
			overlap:  overlap,
		})
	}

	for i := 0; i < b.N; i++ {
		copy(data, m)
		b.StartTimer()

		sort.SliceStable(data, func(i, j int) bool {
			return data[i].overlap > data[j].overlap
		})

		b.StopTimer()
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
