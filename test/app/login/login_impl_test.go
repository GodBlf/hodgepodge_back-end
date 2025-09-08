package login

import (
	"go.uber.org/zap"
	"testing"
	"xmu_roll_call/app/client"
	"xmu_roll_call/app/encrypt"
	"xmu_roll_call/app/login"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"
)

func InitializeLoginImpl() *login.LoginImpl {
	restyClient := client.NewClient()
	encryptImpl := encrypt.NewEncryptImpl()
	loginImpl := login.NewLoginImpl(restyClient, encryptImpl)
	return loginImpl
}
func TestLogin(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	loginImpl := InitializeLoginImpl()
	s, b, err := loginImpl.Login(global.Config.UserName, global.Config.PassWord)
	if err != nil {
		zap.L().Error("login err", zap.Error(err))
		return
	}
	zap.L().Info("login success", zap.String("salt", s), zap.Bool("bool", b))
}
