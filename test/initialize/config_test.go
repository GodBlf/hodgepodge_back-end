package initialize

import (
	"testing"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"
	"xmu_roll_call/utils"

	"go.uber.org/zap"
)

func TestConfig(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	tmp := global.Config.UserName
	zap.L().Info(tmp)
}
func TestConfigLlm(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	model, str, err := utils.RandomLlmModel()
	if err != nil {
		zap.L().Error("RandomLlmModel err", zap.Error(err))
		return
	}
	zap.L().Info("RandomLlmModel", zap.String("api_key", model.ApiKey),
		zap.String("url", model.Url),
		zap.String("model", str),
	)
}
