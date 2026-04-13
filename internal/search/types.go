package search

type searchTask struct {
	docHash string
	chunkID int
	vec     []float32
}

type Result struct {
	DocHash string
	ChunkID int
	Score   float64
}

// Doc(hash) -> Chunk(id) -> Embedding
// This structure allows us to perform
// search across multiple documents.
type index map[string]map[int][]float32
