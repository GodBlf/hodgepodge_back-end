package chatsession

import "context"

type ChatSessionInterface interface {
	Chat(ctx context.Context, userMsg string) (string, error)
	//直接把error发货到管道里
	ChatStream(ctx context.Context, userMsg string)
	ChatStreamOut()
	SetSystemPrompt(prompt string)
	ClearHistory()
}

// model
type ChatStreamElement struct {
	Content string
	Err     error
}
