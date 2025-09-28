package eino

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"strings"
	"sync"
	"testing"
	"xmu_roll_call/app_plus/einochat"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

func TestEino(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	chatSession, err := einochat.NewChatSessionImpl()
	if err != nil {
		zap.L().Error("failed to create chat session", zap.Error(err))
		return
	}
	ctx := context.Background()
	chatSession.SetPrompt(ctx, "你是一个{{.job}}老师", "请简要回答以下问题：{.question}",
		map[string]any{"job": "计算机科学"},
	)
	//chatSession.History = append(chatSession.History, schema.UserMessage("你好"))
	//stream, err := chatSession.Llm.Stream(ctx, chatSession.History)
	//for true {
	//	response, err := stream.Recv()
	//	if err != nil {
	//		zap.L().Error("Stream Recv error", zap.Error(err))
	//		break
	//	}
	//	if response == nil {
	//		break
	//	}
	//	zap.L().Info("Stream Recv", zap.String("content", response.Content))
	//}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		chatSession.ChatStream(ctx, "全面讲解gomock库")

	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		chatSession.ChatStreamOut()
	}()
	wg.Wait()
	chatSession.SetPrompt(ctx, "你是一个{{.job}}老师", "请简要回答以下问题：{.question}",
		map[string]any{"job": "计算机科学"},
	)
	reply, err := chatSession.Chat(ctx, "你是谁？")
	if err != nil {
		zap.L().Error("Chat error", zap.Error(err))
		return
	}
	zap.L().Info("Chat reply", zap.String("reply", reply))
	chat, err := chatSession.Chat(ctx, "请记住我叫许昊龙")
	if err != nil {
		zap.L().Error("Chat error", zap.Error(err))
		return
	}
	zap.L().Info("Chat reply", zap.String("reply", chat))
	chat2, err := chatSession.Chat(ctx, "我是谁？")
	if err != nil {
		zap.L().Error("Chat error", zap.Error(err))
		return
	}
	zap.L().Info("Chat reply", zap.String("reply", chat2))

}

func TestEino2(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	tem := float32(0.8)
	config := &openai.ChatModelConfig{
		APIKey:      global.LlmApiKey,
		BaseURL:     global.LlmUrl + "/" + global.OpenaiVersion,
		Model:       global.LlmModel,
		Temperature: &tem,
	}
	ctx := context.Background()
	llm, err := openai.NewChatModel(ctx, config)
	if err != nil {
		zap.L().Error("failed to create chat model", zap.Error(err))
		return
	}
	input := []*schema.Message{
		schema.SystemMessage("你是一个大学老师"),
		schema.UserMessage("介绍下go语言"),
		schema.AssistantMessage(`
go语言是谷歌公司开发的一种静态强类型、编译式、并发性强且具有垃圾回收功能的编程语言。
它的设计目标是简洁、高效和易于维护，适用于构建高性能的服务器端应用程序。`, nil),
		schema.UserMessage("介绍下python语言"),
	}
	stream, err := llm.Stream(ctx, input)
	if err != nil {
		zap.L().Error("failed to create chat model", zap.Error(err))
		return
	}
	sb := &strings.Builder{}
	for true {
		element, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			zap.L().Error("Stream Recv error", zap.Error(err))
			break
		}
		sb.WriteString(element.Content)
		fmt.Print(element.Content)
	}
	input = append(input, schema.AssistantMessage(sb.String(), nil))
	input = []*schema.Message{}
	tmpl := "你是{{.name}}"
	must := template.Must(template.New("msg").Parse(tmpl))
	bf := &bytes.Buffer{}
	must.Execute(bf, map[string]any{
		"name": "大学老师",
	})
	fmt.Println(bf.String())
}

func TestTemplate(t *testing.T) {
	tmpl := "你是{{.name}}"
	parse, err := template.New("my_template").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	//buffer字节缓冲区相当于strings.builder容器
	bf := &bytes.Buffer{}
	parse.Execute(bf, map[string]any{
		"name": "大学老师",
	})
	fmt.Println(bf.String())
	m := map[string]any{}
	v, ok := m["name"]
	if !ok {
		fmt.Println("key not found")
	} else {
		fmt.Println(v)
	}

}
