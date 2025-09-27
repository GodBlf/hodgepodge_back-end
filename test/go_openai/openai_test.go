package go_openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"
	"xmu_roll_call/app_plus/chatsession"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

//global.LlmApiKey, global.LlmUrl global.LlmModel都是写在yaml配置文件里的
//通过initialize.InitConfig()初始化到global包里

func TestOpenai(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	client := NewThirdPartyClient(global.LlmApiKey, global.LlmUrl, "v1")
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: global.LlmModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "讲解算法中的二分答案法",
				},
			},
		},
	)
	defer stream.Close()
	if err != nil {
		zap.L().Error("ChatCompletion error", zap.Error(err))
		return
	}
	if err != nil {
		zap.L().Error("Stream Recv error", zap.Error(err))
		return
	}
	for true {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			break
		}
		fmt.Print(response.Choices[0].Delta.Content)
	}
}

func NewThirdPartyClient(apiKey string, baseDomain string, version string) *openai.Client {
	// 拼接完整的 BaseURL，格式为 "https://api.example.com/v1"
	//.yaml配置文件里直接写域名即可,版本自己添加v1,v2,后边的路径参数go-openai库自动添加
	baseURL := baseDomain + "/" + version
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = baseURL
	return openai.NewClientWithConfig(cfg)
}

func TestMyopenai(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	chatsession := chatsession.NewChatSessionImpl()
	chatsession.SetSystemPrompt("我在学习golang的库go-openai,你是一个资深的golang开发者,请你帮助我理解go-openai的用法")
	ctx := context.Background()
	reply, err := chatsession.Chat(ctx, "你是谁")
	if err != nil {
		zap.L().Error("Chat error", zap.Error(err))
		return
	}
	zap.L().Info("reply", zap.String("reply", reply))
	chat, _ := chatsession.Chat(ctx, "请记住我的名字我叫许昊龙")
	zap.L().Info("reply", zap.String("reply", chat))
	chat2, _ := chatsession.Chat(ctx, "我是谁")
	zap.L().Info("reply", zap.String("reply", chat2))
	//clear
	chatsession.ClearHistory()
	str := "介绍下chatgpt"
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		chatsession.ChatStream(ctx, str)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		chatsession.ChatStreamOut()
	}()
	wg.Wait()
}
