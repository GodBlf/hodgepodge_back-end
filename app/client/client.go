package client

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http/cookiejar"
	"sync"
	"time"
	"xmu_roll_call/global"
)

var (
	ClientVar *resty.Client
	once      sync.Once
)

func NewClient() *resty.Client {
	once.Do(func() {
		jar, err := cookiejar.New(nil)
		if err != nil {
			zap.L().Error("cookiejar err", zap.Error(err))
		}
		ClientVar = resty.New().
			SetCookieJar(jar).
			// 缺省超时和重试策略可按需设置
			SetTimeout(30*time.Second).
			SetHeader("User-Agent", global.Config.UserAgent)
	})
	return ClientVar
}
