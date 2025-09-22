package client

import (
	"net/http/cookiejar"
	"sync"
	"time"
	"xmu_roll_call/global"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	TestClient *resty.Client
	tconce     = &sync.Once{}
)

func NewTestClient() *resty.Client {
	tconce.Do(func() {
		jar, err := cookiejar.New(nil)
		if err != nil {
			zap.L().Error("cookiejar err", zap.Error(err))
		}
		TestClient = resty.New().
			SetCookieJar(jar).
			SetHeader("User-Agent", global.Config.UserAgent).
			SetTimeout(60 * time.Second)
	})
	return TestClient
}
