package config

const (
	defaultCacheDirPrefix    = "filechat"
	defaultOllamaEndpoint    = "http://localhost:11434"
	defaultEmbeddingLength   = 768
	defaultChunkTokenLength  = 256
	defaultChunkTokenOverlap = 64
	defaultSearchKeepNChunks = 3
	defaultSearchGrowNChunks = 3
	defaultRagPromptTemplate = `
Answer the question using the context below.

Context:
{{ context }}

Question:
{{ question }}
`
)

type Config struct {
	CacheDirPrefix    string
	OllamaEndpoint    string
	ReasoningModel    string
	EmbeddingModel    string
	EmbeddingLength   int
	ChunkTokenLength  int
	ChunkTokenOverlap int
	SearchKeepNChunks int
	RagPromptTemplate string
}

func Init() *Config {
	return &Config{
		CacheDirPrefix:    defaultCacheDirPrefix,
		OllamaEndpoint:    defaultOllamaEndpoint,
		EmbeddingLength:   defaultEmbeddingLength,
		ChunkTokenLength:  defaultChunkTokenLength,
		ChunkTokenOverlap: defaultChunkTokenOverlap,
		SearchKeepNChunks: defaultSearchKeepNChunks,
		RagPromptTemplate: defaultRagPromptTemplate,
	}
}
