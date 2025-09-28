package einochat

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"strings"
	"xmu_roll_call/global"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

// ChatStreamElement 用于流式元素输出
type ChatStreamElement struct {
	Content string
	Err     error
}

// ChatSessionImpl 使用 Eino 封装 LLM
type ChatSessionImpl struct {
	History          []*schema.Message
	systemTemplate   string
	questionTemplate string
	Llm              *openai.ChatModel
	streamChan       chan ChatStreamElement
}

// NewChatSessionImpl 创建一个新的聊天会话
func NewChatSessionImpl() (*ChatSessionImpl, error) {
	config := &openai.ChatModelConfig{
		APIKey:  global.LlmApiKey,
		BaseURL: global.LlmUrl + "/" + global.OpenaiVersion,
		Model:   global.LlmModel,
	}
	model, err := openai.NewChatModel(context.Background(), config)
	if err != nil {
		zap.L().Error("failed to create chat model", zap.Error(err))
		return nil, err
	}
	return &ChatSessionImpl{
		History:    make([]*schema.Message, 0),
		streamChan: make(chan ChatStreamElement, 16),
		Llm:        model,
	}, nil
}

// SetPromptTemplate 设置系统模板和用户问题模板
func (cs *ChatSessionImpl) SetPromptTemplate(system, question string) {
	cs.systemTemplate = system
	cs.questionTemplate = question
}

// ClearHistory 清除历史
func (cs *ChatSessionImpl) ClearHistory() {
	cs.History = make([]*schema.Message, 0)
}

// formatPrompt 根据模板和历史生成消息列表
func (cs *ChatSessionImpl) SetPrompt(ctx context.Context, system string, question string, m map[string]any) error {
	cs.ClearHistory()
	cs.SetPromptTemplate(system, question)
	parse, err := template.New("system").Parse(cs.systemTemplate)
	if err != nil {
		zap.L().Error("failed to parse system template", zap.Error(err))
		return err
	}
	bf := &bytes.Buffer{}
	err = parse.Execute(bf, m)
	if err != nil {
		zap.L().Error("failed to execute system template", zap.Error(err))
		return err
	}
	systemPrompt := bf.String()
	cs.History = append(cs.History, schema.SystemMessage(systemPrompt))
	return nil
}

// todo:没用模板
// Chat 非流式聊天
func (cs *ChatSessionImpl) Chat(ctx context.Context, question string) (string, error) {
	cs.History = append(cs.History, schema.UserMessage(question))
	response, err := cs.Llm.Generate(ctx, cs.History)
	if err != nil {
		zap.L().Error("Llm generate error", zap.Error(err))
		return "", err
	}
	if response.Content == "" {
		return "", errors.New("no response from model")

	}
	reply := response.Content
	cs.History = append(cs.History, schema.AssistantMessage(reply, nil))
	return reply, nil
}

// ChatStream 流式聊天 （发送到 channel）
func (cs *ChatSessionImpl) ChatStream(ctx context.Context, question string) {
	defer close(cs.streamChan)
	cs.History = append(cs.History, schema.UserMessage(question))
	stream, err := cs.Llm.Stream(ctx, cs.History)
	defer stream.Close()
	if err != nil {
		cs.streamChan <- ChatStreamElement{Err: err}
		zap.L().Error("request stream", zap.Error(err))
		return
	}
	reply := &strings.Builder{}
	for true {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			cs.streamChan <- ChatStreamElement{Err: err}
			return
		}
		content := response.Content
		reply.WriteString(content)
		if content != "" {
			cs.streamChan <- ChatStreamElement{Content: content}
		}
	}
	cs.History = append(cs.History, schema.AssistantMessage(reply.String(), nil))

}

// ChatStreamOut 输出流式内容
func (cs *ChatSessionImpl) ChatStreamOut() {
	for elem := range cs.streamChan {
		if elem.Err != nil {
			zap.L().Error("stream error", zap.Error(elem.Err))
			return
		}
		fmt.Print(elem.Content)
	}
}
