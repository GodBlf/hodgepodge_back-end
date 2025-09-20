package client

import (
	"sync"
	"time"
	"xmu_roll_call/global"

	"github.com/go-resty/resty/v2"
)

var (
	LlmClient *resty.Client
	onceLlm   = &sync.Once{}
)

func NewLlmClient() *resty.Client {
	onceLlm.Do(func() {
		LlmClient = resty.New().
			SetTimeout(60*time.Second).
			SetHeader("User-Agent", global.Config.UserAgent)
	})
	return LlmClient
}
