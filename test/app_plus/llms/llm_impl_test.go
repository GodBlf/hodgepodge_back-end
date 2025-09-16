package llms

import (
	"go.uber.org/zap"
	"testing"
	"xmu_roll_call/app_plus/llms"
	"xmu_roll_call/initialize"
)

func TestTr(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitRouter()
	oput, err := llms.Translate("你好")
	if err != nil {
		zap.L().Error(
			"Translate failed: ", zap.Error(err),
		)
		return
	}
	zap.L().Info("Translate response: " + oput)
}
