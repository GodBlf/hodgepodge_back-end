package roll_call

type RollCall interface {
	RollCallStatus() (map[string]int, error)
}
