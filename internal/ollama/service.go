package ollama

import (
	"context"
	"net/http"
	"net/url"

	"filechat/internal/config"

	"github.com/ollama/ollama/api"
)

type Service struct {
	api  *api.Client
	conf *config.Config
}

func New(
	ctx context.Context,
	conf *config.Config,
) (*Service, error) {
	endpoint := conf.OllamaEndpoint
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	ollama := api.NewClient(url, client)
	err = ollama.Heartbeat(ctx)
	if err != nil {
		return nil, err
	}
	return &Service{
		api:  ollama,
		conf: conf,
	}, nil
}

func (s *Service) Embed(
	ctx context.Context,
	chunks []string,
) ([][]float32, error) {
	req := &api.EmbedRequest{
		Model:      s.conf.EmbeddingModel,
		Dimensions: s.conf.EmbeddingLength,
		Input:      chunks,
	}
	res, err := s.api.Embed(ctx, req)
	if err != nil {
		return nil, err
	}
	vecs := res.Embeddings
	return vecs, nil
}

func (s *Service) Chat(
	ctx context.Context,
	msgs []ChatMessage,
) <-chan ChatReply {
	out := make(chan ChatReply)
	go func() {
		defer close(out)
		fn := func(res api.ChatResponse) error {
			text := res.Message.Content
			reply := ChatReply{Text: text}
			select {
			case <-ctx.Done():
			case out <- reply:
			}
			return nil
		}
		apiMsgs := make([]api.Message, len(msgs))
		for _, msg := range msgs {
			apiMsgs = append(apiMsgs, api.Message{
				Role:    string(msg.Role),
				Content: msg.Text,
			})
		}
		req := &api.ChatRequest{
			Model:    s.conf.ReasoningModel,
			Messages: apiMsgs,
		}
		err := s.api.Chat(ctx, req, fn)
		if err != nil {
			reply := ChatReply{Err: err}
			select {
			case <-ctx.Done():
			case out <- reply:
			}
		}
	}()
	return out
}
