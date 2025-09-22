package main

import (
	"xmu_roll_call/app_plus/llms"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志、配置、路由
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitRouter()

	// 初始化 RollCall 实现
	rollCallImpl := InitializeRollCallImpl()

	// 登录并获取 deviceId
	deviceId, err := rollCallImpl.RollCallLogin()
	if err != nil {
		zap.L().Error("RollCallLogin 失败", zap.Error(err))
		return
	}
	rollCallImpl.DeviceId = deviceId
	zap.L().Info("登录成功", zap.String("deviceId", deviceId))

	// ========== 路由定义 ==========
	// 根路径重定向到 /test/:input
	//global.Router.GET("/:input", func(c *gin.Context) {
	//	input := c.Param("input")
	//	c.Redirect(302, "/testllm/"+input)
	//})

	// 测试 LLM 翻译接口
	global.Router.GET("/testllm/:input", func(c *gin.Context) {
		input := c.Param("input")
		translate, err := llms.Translate(input)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			zap.L().Error("Translate failed", zap.Error(err))
			return
		}
		c.JSON(200, gin.H{
			"response": translate,
		})
	})

	// ping 测试
	global.Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"message": "pong",
		})
	})

	// llm 对话接口
	llmImpl := llms.NewLlmImpl()
	global.Router.POST("/llm", llmImpl.Send)

	//
	// 自动签到（新版 RollCallFinal 支持数字+雷达）
	global.Router.GET("/r", rollCallImpl.RollCallFinal)

	// 启动服务
	zap.L().Info("服务启动，监听 8080 端口")
	if err := global.Router.Run(":8080"); err != nil {
		zap.L().Fatal("服务器启动失败", zap.Error(err))
	}

}
