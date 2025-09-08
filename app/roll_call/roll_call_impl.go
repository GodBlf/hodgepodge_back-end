package roll_call

import (
	"sync"
	"xmu_roll_call/app/login"
)

var (
	RollCallImplVar *RollCallImpl
	once            sync.Once
)

type RollCallImpl struct {
	Login login.Login
}

func NewRollCallImpl(login login.Login) *RollCallImpl {
	once.Do(func() {
		RollCallImplVar = &RollCallImpl{
			Login: login,
		}
	})
	return RollCallImplVar
}
