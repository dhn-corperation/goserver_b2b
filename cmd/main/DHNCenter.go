package main

import (
	"os"
	"fmt"
	
	"mycs/init/systemd"
	config "mycs/configs"
	databasepool "mycs/configs/databasepool"

	"github.com/kardianos/service"
)

const (
	name        = "DHNCenter"
	description = "대형네트웍스 카카오 Center API"
)

func main() {

	config.InitCenterConfig()

	databasepool.InitDatabase(true)

	svcConfig := &service.Config{
		Name:        name,
		DisplayName: name,
		Description: description,
	}

	prg := &systemd.Program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		config.Stdlog.Println("Error: ", err)
		os.Exit(1)
	}

	status, err := systemd.Manage(s)
	if err != nil {
		config.Stdlog.Println(status, " Error: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}