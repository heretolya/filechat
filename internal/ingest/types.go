package ingest

type Chunk struct {
	ID   int       `json:"id"`
	Text string    `json:"text"`
	Vec  []float32 `json:"vec"`
}

type Document struct {
	Hash   string  `json:"hash"`
	Chunks []Chunk `json:"chunks"`
}

type Config struct {
	EmbeddingModel    string `json:"embedding_model"`
	EmbeddingLength   int    `json:"embedding_length"`
	ChunkTokenLength  int    `json:"chunk_token_length"`
	ChunkTokenOverlap int    `json:"chunk_token_overlap"`
}

type Snapshot struct {
	Config   Config   `json:"config"`
	Document Document `json:"document"`
}

type Registry map[string]Snapshot
