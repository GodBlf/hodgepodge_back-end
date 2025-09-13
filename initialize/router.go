package initialize

import (
	"github.com/gin-gonic/gin"
	"xmu_roll_call/global"
)

func InitRouter() {
	global.Router = gin.Default()
}
