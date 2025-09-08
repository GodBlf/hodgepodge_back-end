package initialize

import (
	"go.uber.org/zap"
	"testing"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"
)

func TestConfig(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	tmp := global.Config.UserName
	zap.L().Info(tmp)
}
