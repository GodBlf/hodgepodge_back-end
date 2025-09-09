package roll_call

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"testing"
	"xmu_roll_call/app/client"
	"xmu_roll_call/app/roll_call"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"
	"xmu_roll_call/mocks"
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
