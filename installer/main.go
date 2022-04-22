package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kardianos/service"
)

func checkError(err error) bool {
	if err != nil {
		_, fileName, lineNum, _ := runtime.Caller(1)
		fileName = filepath.Base(fileName)
		fmt.Printf("%s %s:%d: %s \n", time.Now().Format("2006-01-02 15:04:05"), fileName, lineNum, errors.Unwrap(err).Error())
		return true
	}
	return false
}

type program struct{}

func main() {
	svcConfig := &service.Config{
		Name:        "ratd",
		DisplayName: "ratd Agent",
		Description: "IT asset management tool",
		UserName:    "root",
		Executable:  "/usr/local/bin/ratd",
		Option:      service.KeyValue{"KeepAlive": true, "RunAtLoad": true},
	}

	prg := program{}
	s, err := service.New(prg, svcConfig)
	if checkError(err) {
		log.Fatal(err)
	}

	//hc := http.Client{Timeout: time.Minute}
	// temporary to prove install works until we can serve the installer from elsewhere
	agentSrcBinary, err := os.Open("../agent/" + svcConfig.Name)
	if checkError(err) {
		return
	}
	agentDestBinary, err := os.Create(svcConfig.Executable)
	if checkError(err) {
		return
	}

	written, err := io.Copy(agentDestBinary, agentSrcBinary)
	if checkError(err) {
		return
	}
	log.Printf("Wrote %d bytes to '%s'", written, svcConfig.Executable)
	agentSrcBinary.Close()
	agentDestBinary.Close()

	err = os.Chmod(svcConfig.Executable, 755)
	if checkError(err) {
		return
	}

	err = s.Uninstall()
	if checkError(err) {
		//
	}

	log.Printf("installing agent")
	err = s.Install()
	if checkError(err) {
		log.Fatal(err)
	}

	log.Printf("agent installed")

	log.Printf("starting agent")

	err = s.Start()
	if checkError(err) {
		return
	}

	log.Printf("agent started")

}

func (prg program) Stop(s service.Service) error {
	return nil
}

func (prg program) Start(s service.Service) error {
	return nil
}
