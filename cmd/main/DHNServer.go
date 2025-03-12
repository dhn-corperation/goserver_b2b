package main

import (
	"os"
	"fmt"
	"syscall"
	
	"mycs/init/systemd"
	config "mycs/configs"
	databasepool "mycs/configs/databasepool"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kardianos/service"
)

const (
	name        = "DHNServer"
	description = "대형네트웍스 카카오 발송 서버"
)

func main() {

	config.InitConfig()

	databasepool.InitDatabase(false)

	var rLimit syscall.Rlimit

	rLimit.Max = 50000
	rLimit.Cur = 50000

	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		config.Stdlog.Println("Error Setting Rlimit ", err)
	}

	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		config.Stdlog.Println("Error Getting Rlimit ", err)
	}

	config.Stdlog.Println("Rlimit Final", rLimit)

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


