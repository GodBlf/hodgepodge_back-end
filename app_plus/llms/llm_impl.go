package llms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	client2 "xmu_roll_call/app/client"
	"xmu_roll_call/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ChatMessage 表示一条消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 表示 GPT API 请求体
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream,omitempty"`
}

// ChatResponse 表示 GPT API 响应体（只保留必要字段）
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var (
	LlmImplVar *LlmImpl
	once       = &sync.Once{}
)

type LlmImpl struct {
}

func NewLlmImpl() *LlmImpl {
	once.Do(func() {
		LlmImplVar = &LlmImpl{}
	})
	return LlmImplVar
}
func (l *LlmImpl) Send(context *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}
	context.ShouldBindJSON(&req)
	reply, err := Translate(req.Message)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, gin.H{"reply": reply})
}

// Translate 将输入文本发送到 GPT API 并返回翻译结果
func Translate(input string) (string, error) {
	model, modelname, err2 := utils.RandomLlmModel()
	if err2 != nil {
		zap.L().Error("RandomLlmModel failed: ", zap.Error(err2))
		return "", err2
	}
	// 构造请求体
	reqBody := ChatRequest{
		Model: modelname, // 你可以改成实际模型
		Messages: []ChatMessage{
			{Role: "system", Content: "你是一个ai助手"},
			{Role: "user", Content: input},
		},
	}

	client := client2.NewLlmClient()

	request := client.R()
	request.SetBody(reqBody).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+model.ApiKey)
	response, err := request.Post(model.Url)
	if err != nil {
		zap.L().Error(
			"LLM API request error", zap.Error(err),
		)
		return "", err
	}
	if response.StatusCode() != http.StatusOK {
		zap.L().Error(
			"LLM API response error", zap.Int("status", response.StatusCode()),
			zap.String("body", response.String()),
		)
		return "", fmt.Errorf("LLM API error: %s", response.String())
	}

	var chatResp ChatResponse
	json.NewDecoder(bytes.NewBuffer(response.Body())).Decode(&chatResp)

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from GPT API")
	}
	zap.L().Info(
		"LLM response",
		zap.String("reply", chatResp.Choices[0].Message.Content),
	)
	return chatResp.Choices[0].Message.Content, nil
}
