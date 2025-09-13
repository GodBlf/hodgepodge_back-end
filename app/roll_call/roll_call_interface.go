package roll_call

import "github.com/gin-gonic/gin"

type RollCall interface {
	RollCallFinal(c *gin.Context)
	RollCallLogin() (string, error)
	RollCallStatus() (map[string]int, error)
	NumberCodeQuery(rollcall map[string]int) (map[string]string, error, int)
	NumberCodePost(courseNameRollCallId map[string]int, numberCode map[string]string, deviceId string) error
	//todo:radar rollcall
}
