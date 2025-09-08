//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"xmu_roll_call/app/client"
	"xmu_roll_call/app/encrypt"
	"xmu_roll_call/app/login"
)

func InitializeLoginImpl() *login.LoginImpl {
	wire.Build(client.NewClient, encrypt.NewEncryptImpl, login.NewLoginImpl,
		wire.Bind(new(encrypt.Encrypt), new(*encrypt.EncryptImpl)),
	)
	return &login.LoginImpl{}

}
