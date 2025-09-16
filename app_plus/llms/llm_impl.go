package llms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sync"
	client2 "xmu_roll_call/app/client"

	"xmu_roll_call/global"
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
	// 构造请求体
	reqBody := ChatRequest{
		Model: global.LlmModel, // 你可以改成实际模型
		Messages: []ChatMessage{
			{Role: "system", Content: "你是一个ai助手"},
			{Role: "user", Content: input},
		},
	}

	//bodyBytes, err := json.Marshal(reqBody)
	//if err != nil {
	//	return "", err
	//}

	client := client2.NewLlmClient()

	request := client.R()
	request.SetBody(reqBody).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+global.LlmApiKey)
	//req, err := http.NewRequest("POST", global.LlmUrl, bytes.NewBuffer(bodyBytes))
	//if err != nil {
	//	return "", err
	//}

	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer "+global.LlmApiKey)

	//resp, err := client.Do(req)
	//if err != nil {
	//	return "", err
	//}
	//defer resp.Body.Close()
	response, err := request.Post(global.LlmUrl)
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
	//if resp.StatusCode != http.StatusOK {
	//	b, _ := io.ReadAll(resp.Body)
	//	return "", fmt.Errorf("GPT API error: %s", string(b))
	//}

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
