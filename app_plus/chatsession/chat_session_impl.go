package chatsession

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"xmu_roll_call/global"

	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// 需要lazyglobal时候实现单例
var (
	ChatSessionImplVar *ChatSessionImpl = &ChatSessionImpl{}
	csionce                             = &sync.Once{}
)

type ChatSessionImpl struct {
	Messages       []openai.ChatCompletionMessage
	client         *openai.Client
	model          string
	ChatStreamChan chan ChatStreamElement
}

func NewChatSessionImpl() *ChatSessionImpl {
	config := openai.DefaultConfig(global.LlmApiKey)
	config.BaseURL = global.LlmUrl + "/" + global.OpenaiVersion
	return &ChatSessionImpl{
		Messages:       make([]openai.ChatCompletionMessage, 0),
		client:         openai.NewClientWithConfig(config),
		model:          global.LlmModel,
		ChatStreamChan: make(chan ChatStreamElement, 2),
	}
}

func (cs *ChatSessionImpl) Chat(ctx context.Context, userMsg string) (string, error) {
	cs.AddMessage(openai.ChatMessageRoleUser, userMsg)
	response, err := cs.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    cs.model,
		Messages: cs.Messages,
	})
	if err != nil {
		zap.L().Error("ChatCompletion error", zap.Error(err))
		return "", err
	}
	if len(response.Choices) == 0 {
		return "", errors.New("no response from model")
	}
	reply := response.Choices[0].Message.Content
	cs.AddMessage(openai.ChatMessageRoleAssistant, reply)
	return reply, nil
}

func (cs *ChatSessionImpl) ChatStream(ctx context.Context, userMsg string) {
	cs.AddMessage(openai.ChatMessageRoleUser, userMsg)
	stream, err := cs.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    cs.model,
		Messages: cs.Messages,
	})
	if err != nil {
		cs.ChatStreamChan <- ChatStreamElement{Err: err}
		close(cs.ChatStreamChan)
		zap.L().Error("request stream", zap.Error(err))
		return
	}
	defer stream.Close()

	reply := &strings.Builder{}
	for true {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			cs.ChatStreamChan <- ChatStreamElement{Err: err}
			close(cs.ChatStreamChan)
			return
		}
		delta := response.Choices[0].Delta.Content
		if delta != "" {
			reply.WriteString(delta)
			cs.ChatStreamChan <- ChatStreamElement{Content: delta}
		}
	}
	cs.AddMessage(openai.ChatMessageRoleAssistant, reply.String())
	close(cs.ChatStreamChan)
}

func (cs *ChatSessionImpl) ChatStreamOut() {
	for element := range cs.ChatStreamChan {
		if element.Err != nil {
			zap.L().Error("stream error", zap.Error(element.Err))
			return
		}
		fmt.Print(element.Content)
	}
}

func (cs *ChatSessionImpl) AddMessage(role, content string) {
	cs.Messages = append(cs.Messages, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
}

func (cs *ChatSessionImpl) ClearHistory() {
	cs.Messages = make([]openai.ChatCompletionMessage, 0)
}
func (cs *ChatSessionImpl) SetSystemPrompt(prompt string) {
	cs.ClearHistory()
	cs.AddMessage(openai.ChatMessageRoleSystem, prompt)
}
