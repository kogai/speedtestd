package main

import (
	"fmt"
	"github.com/kogai/speedtestd"
	"github.com/takama/daemon"
	"log"
	"os"
)

const (
	name        = "speedtestd"
	description = "speedtest-go as Service."
)

var dependencies = []string{"dummy.service"}
var stdlog, errlog *log.Logger

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	rawService, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &speedtestd.Service{rawService}
	// service.port = port

	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
