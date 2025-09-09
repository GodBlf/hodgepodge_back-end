package main

import (
	"go.uber.org/zap"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	loginImpl := InitializeLoginImpl()
	rollCallImpl := InitializeRollCallImpl()
	_, b, err := loginImpl.Login(global.Config.UserName, global.Config.PassWord)
	if err != nil || !b {
		return
	}
	m, err := rollCallImpl.RollCallStatus()

}
