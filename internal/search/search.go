package search

import (
	"context"
	"runtime"
	"sort"
	"sync"

	"filechat/internal/tools"

	"github.com/charmbracelet/log"
)

type Service struct {
	mu    *sync.RWMutex
	index index
	total int
}

func New() *Service {
	return &Service{
		mu:    &sync.RWMutex{},
		index: make(index),
		total: 0,
	}
}

func (s *Service) Add(
	docHash string,
	chunkID int,
	vec []float32,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.index[docHash]; !ok {
		chunk := make(map[int][]float32)
		s.index[docHash] = chunk
	}
	s.index[docHash][chunkID] = vec
	s.total++
}

func (s *Service) Search(
	ctx context.Context,
	q []float32, n int,
) ([]Result, error) {
	results := make([]Result, 0)
	for r := range s.pSearch(ctx, q) {
		results = append(results, r)
	}
	// Best scores will appear first
	sort.Slice(results, func(i, j int) bool {
		a := results[i].Score
		b := results[j].Score
		return a > b
	})
	out := results[:n]
	return out, nil
}

func (s *Service) traverse(
	ctx context.Context,
) <-chan searchTask {
	out := make(chan searchTask)
	go func() {
		s.mu.RLock()
		defer s.mu.RUnlock()
		defer close(out)
		for docHash, chunks := range s.index {
			for chunkID, vec := range chunks {
				task := searchTask{
					vec:     vec,
					docHash: docHash,
					chunkID: chunkID,
				}
				select {
				case out <- task:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

func (s *Service) pSearch(
	ctx context.Context,
	q []float32,
) <-chan Result {
	var wg sync.WaitGroup
	tasks := s.traverse(ctx)
	out := make(chan Result, 10)
	for range runtime.NumCPU() {
		wg.Go(func() {
			for t := range tasks {
				if ctx.Err() != nil {
					return
				}
				score, err := tools.CosineSim(q, t.vec)
				if err != nil {
					log.Warn(err)
					continue
				}
				result := Result{
					DocHash: t.docHash,
					ChunkID: t.chunkID,
					Score:   score,
				}
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
