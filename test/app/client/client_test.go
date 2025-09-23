package client

import (
	"net/http/cookiejar"
	"sync"
	"testing"
	"time"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	TestClient *resty.Client
	tconce     = &sync.Once{}
)

func NewTestClient() *resty.Client {
	initialize.InitConfig()
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

func TestClientlunxun(t *testing.T) {
	initialize.InitLogger()
	client := NewTestClient()
	for true {
		r := client.R()
		response, err := r.Get("http://10.242.194.32:8080/r")
		if err != nil {
			zap.L().Error("请求失败", zap.Error(err))
		}
		zap.L().Info("请求成功", zap.String("响应", response.String()))
		time.Sleep(1 * time.Second)
	}
}
