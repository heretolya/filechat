package chat

import "filechat/internal/ollama"

type messages []ollama.ChatMessage

type userInput struct {
	err  error
	text string
}
