package app

import (
	"sync"
	"xmu_roll_call/app/login"
	"xmu_roll_call/app/roll_call"
)

var (
	AppImplVar *AppImpl
	once       sync.Once
)

type AppImpl struct {
	Login    login.Login
	RollCall roll_call.RollCall
}

func (a *AppImpl) AppLogin() error {
	//TODO implement me
	panic("implement me")
}

func (a *AppImpl) AppImplRollCall() error {
	//TODO implement me
	panic("implement me")
}

func NewAppImpl(login login.Login, rollCall roll_call.RollCall) *AppImpl {
	once.Do(func() {
		AppImplVar = &AppImpl{
			Login:    login,
			RollCall: rollCall,
		}
	})
	return AppImplVar
}
