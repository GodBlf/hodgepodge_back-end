package users_client

import (
	"net/http/cookiejar"
	"sync"
	"time"
	"xmu_roll_call/global"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	Client *resty.Client
	conce  = &sync.Once{}
)

func NewClient() *resty.Client {
	conce.Do(func() {
		jar, err := cookiejar.New(nil)
		if err != nil {
			zap.L().Error("cookiejar err", zap.Error(err))
		}
		header := resty.New().SetCookieJar(jar).
			SetTimeout(60*time.Second).
			SetHeader("User-Agent", global.Config.UserAgent)
		Client = header
	})
	return Client
}

func ClearClientCookiejar() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		zap.L().Error("cookiejar err", zap.Error(err))
	}
	Client.SetCookieJar(jar)
}
