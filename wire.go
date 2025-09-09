//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"xmu_roll_call/app"
	"xmu_roll_call/app/client"
	"xmu_roll_call/app/encrypt"
	"xmu_roll_call/app/login"
	"xmu_roll_call/app/roll_call"
)

func InitializeLoginImpl() *login.LoginImpl {
	wire.Build(client.NewClient, encrypt.NewEncryptImpl, login.NewLoginImpl,
		wire.Bind(new(encrypt.Encrypt), new(*encrypt.EncryptImpl)),
	)
	return &login.LoginImpl{}

}

func InitializeRollCallImpl() *roll_call.RollCallImpl {
	wire.Build(client.NewClient, roll_call.NewRollCallImpl)
	return &roll_call.RollCallImpl{}
}

func InitializeAppImpl() *app.AppImpl {
	wire.Build(login.NewLoginImpl, roll_call.NewRollCallImpl, app.NewAppImpl, client.NewClient, encrypt.NewEncryptImpl,
		wire.Bind(new(encrypt.Encrypt), new(*encrypt.EncryptImpl)),
		wire.Bind(new(login.Login), new(*login.LoginImpl)),
		wire.Bind(new(roll_call.RollCall), new(*roll_call.RollCallImpl)),
	)
	return &app.AppImpl{}
}
