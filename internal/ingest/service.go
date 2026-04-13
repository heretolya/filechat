package ingest

import (
	"context"
	"fmt"
	"strings"

	"filechat/internal/config"
	"filechat/internal/ollama"
	"filechat/internal/tools"
)

type Service struct {
	conf *config.Config
	ai   *ollama.Service
}

func New(
	conf *config.Config,
	ai *ollama.Service,
) *Service {
	return &Service{
		conf: conf,
		ai:   ai,
	}
}

func (s *Service) Hash(text string) string {
	var buf strings.Builder
	buf.WriteString(text)
	buf.WriteString(s.conf.EmbeddingModel)
	fmt.Fprint(&buf, fmt.Sprint(s.conf.EmbeddingLength))
	fmt.Fprint(&buf, fmt.Sprint(s.conf.ChunkTokenLength))
	fmt.Fprint(&buf, fmt.Sprint(s.conf.ChunkTokenOverlap))
	return tools.Hash(buf.String())
}

func (s *Service) Do(
	ctx context.Context,
	text string,
) (*Snapshot, error) {
	tokenChunks, err := tools.SliceToChunks(
		tools.Tokenize(text),
		s.conf.ChunkTokenLength,
		s.conf.ChunkTokenOverlap,
	)
	if err != nil {
		return nil, err
	}
	chunks := make([]Chunk, 0, len(tokenChunks))
	inputs := make([]string, 0, len(tokenChunks))
	for idx, textChunk := range tokenChunks {
		text := strings.Join(textChunk, " ")
		chunk := Chunk{ID: idx, Text: text}
		chunks = append(chunks, chunk)
		inputs = append(inputs, text)
	}
	vecs, err := s.ai.Embed(ctx, inputs)
	if err != nil {
		return nil, err
	}
	for idx, vec := range vecs {
		chunks[idx].Vec = vec
	}
	return &Snapshot{
		Config: Config{
			EmbeddingModel:    s.conf.EmbeddingModel,
			EmbeddingLength:   s.conf.EmbeddingLength,
			ChunkTokenLength:  s.conf.ChunkTokenLength,
			ChunkTokenOverlap: s.conf.ChunkTokenOverlap,
		},
		Document: Document{
			Hash:   s.Hash(text),
			Chunks: chunks,
		},
	}, nil
}
