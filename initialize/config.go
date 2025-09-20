package initialize

import (
	"xmu_roll_call/global"
	"xmu_roll_call/model"

	"github.com/spf13/viper"
)

type ConfigModel struct {
	UserName  string
	PassWord  string
	UserAgent string
}

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("D:\\tmp_test\\qiandaotest\\config\\config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic("读取配置文件失败: " + err.Error())
	}
	global.LlmUrl = v.GetString("llm_url")
	global.LlmApiKey = v.GetString("llm_api_key")
	global.LlmModel = v.GetString("llm_model")
	username := v.GetString("username")
	password := v.GetString("password")
	ua := v.GetString("user_agent")
	if ua == "" {
		ua = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Mobile Safari/537.36"
	}
	if username == "" || password == "" {
		panic("配置文件缺少 username/password")
	}
	cfg := &model.ConfigModel{
		UserName:          username,
		PassWord:          password,
		UserAgent:         ua,
		IdsUrl:            v.GetString("ids_url"),
		RollCallStatusUrl: v.GetString("roll_call_status_url"),
	}
	//llmlist config
	global.MLUK = &model.ModelLlmUrlKey{
		ModelCount: 2,
	}
	v.Unmarshal(global.MLUK)
	global.Config = cfg
}
