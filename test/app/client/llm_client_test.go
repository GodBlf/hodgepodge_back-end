package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"xmu_roll_call/app/client"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"
	"xmu_roll_call/utils"

	"go.uber.org/zap"
)

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
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func TestRandomLlmFetch(t *testing.T) {
	initialize.InitConfig()
	initialize.InitLogger()
	llmClient := client.NewLlmClient()
loop:
	model, modelname, err := utils.RandomLlmModel()
	if modelname == "claude-4-sonnet" {
		goto loop
	}
	if err != nil {
		zap.L().Error("RandomLlmFetch failed: ", zap.Error(err))
	}
	request := llmClient.R()
	reqBody := ChatRequest{
		Model: modelname, // 你可以改成实际模型
		Messages: []ChatMessage{
			{Role: "system", Content: "你是一个ai助手"},
			{Role: "user", Content: "你好"},
		},
	}
	request.SetBody(reqBody).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+model.ApiKey)
	response, err := request.Post(global.LlmUrl)
	if err != nil {
		zap.L().Error(
			"LLM API request error", zap.Error(err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		zap.L().Error(
			"LLM API response error", zap.Int("status", response.StatusCode()),
			zap.String("url", model.Url),
			zap.String("model", modelname),
			zap.String("api_key", model.ApiKey),
			zap.String("body", response.String()),
		)
		return
	}
	var chatResp ChatResponse
	json.NewDecoder(bytes.NewBuffer(response.Body())).Decode(&chatResp)

	if len(chatResp.Choices) == 0 {
		zap.L().Error("LLM API response has no choices")
		return
	}
	zap.L().Info(
		"LLM response",
		zap.String("reply", chatResp.Choices[0].Message.Content),
	)
	return
}
