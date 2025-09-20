package main

import (
	"fmt"
	"time"
	"xmu_roll_call/app_plus/llms"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitRouter()
	rollCallImpl := InitializeRollCallImpl()
	deviceId, err := rollCallImpl.RollCallLogin()
	rollCallImpl.DeviceId = deviceId
	if err != nil {
		zap.L().Error("RollCallLogin err", zap.Error(err))
		return
	}
	zap.L().Info("登录成功,查询签到状态...")
	//"/"重定向为"/test" get方法
	global.Router.GET("/:input", func(c *gin.Context) {
		input := c.Param("input")
		c.Redirect(302, "/test/"+input)
	})
	global.Router.GET("/test/:input", func(c *gin.Context) {
		input := c.Param("input")
		translate, err2 := llms.Translate(input)
		if err2 != nil {
			c.JSON(500, gin.H{"error": err2.Error()})
			zap.L().Error("Translate failed", zap.Error(err2))
		}
		c.JSON(200, gin.H{
			"response": translate,
		})
	})
	//ping
	global.Router.GET("/ping", func(context *gin.Context) {
		context.JSON(200, map[string]string{
			"message": "pong",
		})
	})
	//llm
	llmImpl := llms.NewLlmImpl()
	global.Router.POST("/llm", llmImpl.Send)
	//
	global.Router.GET("/r", rollCallImpl.RollCallFinal)
	global.Router.Run(":8080")
	return
	//todo:delete
	if err != nil {
		zap.L().Error("RollCallLogin err", zap.Error(err))
		return
	}
	zap.L().Info("登录成功,查询签到状态...")
	for {
		courseAndRollcallId, err := rollCallImpl.RollCallStatus()
		if err != nil {
			zap.L().Error("查询签到状态失败", zap.Error(err))
			return
		}
		zap.L().Info("查询签到状态成功", zap.Any("查询结果", courseAndRollcallId))
		if len(courseAndRollcallId) == 0 {
			zap.L().Info("当前没有需要签到的课程")
			time.Sleep(2 * time.Second)
			continue
		}
		zap.L().Info("开始查询签到码...")
		query, err, rollcallType := rollCallImpl.NumberCodeQuery(courseAndRollcallId)

		if rollcallType == 0 {
			zap.L().Info("签到码查询成功", zap.Any("查询结果", query))
			zap.L().Info("打印签到码...")
			for title, numbercode := range query {
				if numbercode != "" {
					fmt.Printf("✅ %s: 签到码 %s\n", title, numbercode)
				} else {
					fmt.Printf("❌ %s: 获取签到码失败\n", title)
				}
			}
			zap.L().Info("开启自动签到...")

			err = rollCallImpl.NumberCodePost(courseAndRollcallId, query, deviceId)
			if err != nil {
				zap.L().Error("自动签到失败", zap.Error(err))
				break
			}
			zap.L().Info("自动签到成功")
			//签到完成后跳出循环
			break

		} else if rollcallType == 1 {
			//radar rollcall
			zap.L().Info("雷达签到", zap.Any("查询结果", query))
			zap.L().Info("开启自动签到...")
			time.Sleep(200 * time.Second)
			break
		} else {
			zap.L().Warn("二维码签到请老实上课", zap.Any("查询结果", query))
			break
		}
	}
	time.Sleep(time.Second * 200)
}
