package users_roll_call

import "github.com/gin-gonic/gin"

type UersRollCallInterface interface {
	RollCall(c *gin.Context)
}
