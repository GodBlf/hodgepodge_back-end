package global

import (
	"github.com/gin-gonic/gin"
	"xmu_roll_call/model"
)

var (
	Config *model.ConfigModel
	Router *gin.Engine
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
