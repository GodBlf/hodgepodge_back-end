package roll_call

type RollCall interface {
	RollCallLogin() (string, error)
	RollCallStatus() (map[string]int, error)
	NumberCodeQuery(rollcall map[string]int) (map[string]string, error, int)
	NumberCodePost(courseNameRollCallId map[string]int, numberCode map[string]string, deviceId string) error
	//todo:radar rollcall
}
