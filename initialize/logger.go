package initialize

import (
	"go.uber.org/zap"
)

func InitLogger() {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
