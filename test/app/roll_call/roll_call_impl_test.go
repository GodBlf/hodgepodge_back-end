package roll_call

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"testing"
	"xmu_roll_call/app/client"
	"xmu_roll_call/app/roll_call"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"
	"xmu_roll_call/mocks"
	"xmu_roll_call/model"
)

func TestRollCallLogin(t *testing.T) {
	initialize.InitConfig()
	initialize.InitConfig()
	c := client.NewClient()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogin := mocks.NewMockLogin(ctrl)
	mockLogin.EXPECT().Login(global.Config.UserName, global.Config.PassWord).Return("123_123", true, nil)
	r := roll_call.NewRollCallImpl(c, mockLogin)
	//
	login, err := r.RollCallLogin()
	zap.L().Error("err", zap.Error(err))
	zap.L().Info("login", zap.String("login", login))

}

func TestNumberCodeQuery(t *testing.T) {
	results := make(map[string]string)
	responseString := "{\"is_radar\":false,\"number_code\":\"123456\"}" //response.String()
	isRadar := gjson.Get(responseString, "is_radar").Bool()

	NumberCode := gjson.Get(responseString, "number_code").String()
	fmt.Println(results)
	fmt.Println(isRadar)
	fmt.Println(NumberCode)

}

func TestFinal(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitRouter()
	controller := gomock.NewController(t)
	defer controller.Finish()
	roll := mocks.NewMockRollCall(controller)
	roll.EXPECT().RollCallStatus().Return(
		map[string]int{
			"测试课程": 123,
		}, nil,
	).AnyTimes()
	roll.EXPECT().NumberCodeQuery(
		map[string]int{
			"测试课程": 123,
		},
	).Return(
		map[string]string{
			"测试课程": "123456",
		}, nil, 0,
	).AnyTimes()
	roll.EXPECT().NumberCodePost(nil, nil, "").Return(nil).AnyTimes()
	global.Router.GET("/r", func(c *gin.Context) {
		courseAndRollcallId, err := roll.RollCallStatus()
		if err != nil {
			zap.L().Error("查询签到状态失败", zap.Error(err))
			c.JSON(500, gin.H{"message": "查询签到状态失败"})
			return
		}
		zap.L().Info("查询签到状态成功", zap.Any("查询结果", courseAndRollcallId))
		if len(courseAndRollcallId) == 0 {
			zap.L().Info("当前没有需要签到的课程")
			c.JSON(200, gin.H{"message": "当前没有需要签到的课程"})
			return
		}
		zap.L().Info("开始查询签到码...")
		query, err, rollcallType := roll.NumberCodeQuery(courseAndRollcallId)

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

			err = roll.NumberCodePost(nil, nil, "")
			if err != nil {
				zap.L().Error("自动签到失败", zap.Error(err))
				query["message"] = "自动签到失败"
				c.JSON(500, query)
				return
			}
			zap.L().Info("自动签到成功")
			query["message"] = "自动签到成功"
			c.JSON(200, query)
			return

		} else if rollcallType == 1 {
			//radar rollcall
			zap.L().Info("雷达签到", zap.Any("查询结果", query))
			zap.L().Info("开启自动签到...")
			c.JSON(200, gin.H{"message": "雷达签到"})
			return
		} else {
			zap.L().Warn("二维码签到请老实上课", zap.Any("查询结果", query))
			c.JSON(200, gin.H{"message": "二维码签到请老实上课"})
			return
		}
	})
	global.Router.Run(":8080")
}

func TestStatus(t *testing.T) {
	rollcalls := &model.RollCalls{List: []model.RollCallJson{
		{RollcallID: 123, CourseTitle: "你好"},
	}}
	marshal, _ := json.Marshal(rollcalls)
	s := string(marshal)
	fmt.Println(s)
	unmarshalList := &model.RollCalls{}
	json.Unmarshal(marshal, unmarshalList)
	pending := make(map[string]int)
	for _, status := range unmarshalList.List {
		if status.RollcallID != 0 {
			pending[status.CourseTitle] = status.RollcallID
		}
	}
	for s, i := range pending {
		fmt.Printf("%s: %d\n", s, i)
	}
}
