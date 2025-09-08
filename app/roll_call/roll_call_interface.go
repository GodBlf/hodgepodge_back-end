package roll_call

type RollCall interface {
	RollCallStatus() (map[string]int, error)
	NumberCodeQuery(rollcall map[string]int) (map[string]string, error)
	NumberCodePost(courseNameRollCallId map[string]int, numberCode map[string]string, deviceId string) error
}
