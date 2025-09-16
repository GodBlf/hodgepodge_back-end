package client

import (
	"github.com/go-resty/resty/v2"
	"sync"
	"time"
	"xmu_roll_call/global"
)

var (
	LlmClient *resty.Client
	onceLlm   = &sync.Once{}
)

func NewLlmClient() *resty.Client {
	onceLlm.Do(func() {
		LlmClient = resty.New().
			SetTimeout(30*time.Second).
			SetHeader("User-Agent", global.Config.UserAgent)
	})
	return LlmClient
}
