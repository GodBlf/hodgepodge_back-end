package eino

import (
	"github.com/cloudwego/eino-ext/components/tool/browseruse"
	"go.uber.org/zap"

	"context"
	"testing"
	"xmu_roll_call/initialize"
)

func TestTool(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	url := "https://bilibili.com"
	but, err := browseruse.NewBrowserUseTool(context.Background(), &browseruse.Config{})
	if err != nil {
		zap.L().Error("failed to create browser tool", zap.Error(err))
		return
	}
	execute, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		zap.L().Error("failed to execute browser tool", zap.Error(err))
		return
	}
	zap.L().Info("browser tool execute result", zap.Any("result", execute))

}
