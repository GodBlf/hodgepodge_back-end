package global

import (
	"xmu_roll_call/model"

	"github.com/gin-gonic/gin"
)

var (
	Config *model.ConfigModel
	Router *gin.Engine
)

var (
	LlmUrl    string
	LlmApiKey string
	LlmModel  string
)

var Locations = []struct {
	Latitude  string
	Longitude string
	Name      string
}{
	{"123", "123", "北京大学"},
	{"39.989643", "116.305408", "清华大学"},
	{"31.326315", "121.441049", "上海交通大学"},
}

var MLUK *model.ModelLlmUrlKey
