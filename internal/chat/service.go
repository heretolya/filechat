package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"filechat/internal/config"
	"filechat/internal/ingest"
	"filechat/internal/ollama"
	"filechat/internal/search"

	"github.com/charmbracelet/log"
)

const skipRag = "!!"

type Service struct {
	conf *config.Config
	ai   *ollama.Service
	db   ingest.Registry
	vss  *search.Service
	buf  *strings.Builder
	msgs messages
}

func New(
	conf *config.Config,
	ai *ollama.Service,
	db ingest.Registry,
	vss *search.Service,
) *Service {
	msgs := make(messages, 0)
	buf := &strings.Builder{}
	return &Service{
		ai:   ai,
		db:   db,
		buf:  buf,
		vss:  vss,
		conf: conf,
		msgs: msgs,
	}
}

func (s *Service) input(ctx context.Context) <-chan userInput {
	out := make(chan userInput, 1)
	go func() {
		defer close(out)
		scanner := bufio.NewScanner(os.Stdout)
		for {
			if ctx.Err() != nil {
				return
			}
			if !scanner.Scan() {
				return
			}
			out <- userInput{
				err:  scanner.Err(),
				text: scanner.Text(),
			}
		}
	}()
	return out
}

func (s *Service) rag(ctx context.Context, text string) (string, error) {
	em, err := s.ai.Embed(ctx, []string{text})
	if err != nil {
		return "", err
	}
	n := s.conf.SearchKeepNChunks
	found, err := s.vss.Search(ctx, em[0], n)
	if err != nil {
		return "", err
	}

	// N-best search results are user to
	// enrich prompt with retrieed data,
	// each item should preserve score
	var buf strings.Builder
	for i, item := range found {
		rec, ok := s.db[item.DocHash]
		if !ok {
			continue
		}
		chunks := rec.Document.Chunks
		if item.ChunkID > len(chunks) {
			continue
		}
		chunk := chunks[item.ChunkID]
		templ := "#%d (score=%.2f)\n%s\n"
		score, text := item.Score, chunk.Text
		fmt.Fprintf(&buf, templ, i, score, text)
	}

	prompt := s.conf.RagPromptTemplate
	prompt = strings.ReplaceAll(prompt, "{{ context }}", buf.String())
	prompt = strings.ReplaceAll(prompt, "{{ question }}", text)
	return prompt, nil
}

func (s *Service) ask(ctx context.Context, text string) error {
	s.msgs = append(s.msgs, ollama.ChatMessage{
		Role: ollama.RoleUser,
		Text: text,
	})
	stream := s.ai.Chat(ctx, s.msgs)
	var buf strings.Builder
	for chunk := range stream {
		if chunk.Err != nil {
			return chunk.Err
		}
		buf.WriteString(chunk.Text)
		fmt.Print(chunk.Text)
	}
	result := buf.String()
	s.msgs = append(s.msgs, ollama.ChatMessage{
		Role: ollama.RoleSystem,
		Text: result,
	})
	if !strings.HasSuffix(result, "\n") {
		fmt.Println()
	}
	return nil
}

func (s *Service) Start(ctx context.Context) error {
	fmt.Print("> ")
	input := s.input(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case q, ok := <-input:
			if !ok {
				return nil
			}
			if q.err != nil {
				return q.err
			}
			log.Info("Thinking...")
			text := q.text
			if !strings.HasPrefix(text, skipRag) {
				prompt, err := s.rag(ctx, text)
				if err != nil {
					return err
				}
				text = prompt
			}
			err := s.ask(ctx, text)
			if err != nil {
				return err
			}
			fmt.Print("> ")
		}
	}
}
