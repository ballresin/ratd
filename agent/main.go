package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kardianos/service"
)

type agent struct {
	ID            int
	Version       float64
	hc            http.Client
	checkinServer string
	WANOffline    bool
}

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

	log.SetFlags(0)
	// check if interactive?
	// daemonize yourself
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

	log.Printf("Starting program")

	err = s.Run()
	if checkError(err) {
		return
	}
}

func (p program) Stop(s service.Service) error {
	return nil
}

func (p program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p program) run() {

	var a *agent

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	a = &agent{
		ID:      1,
		Version: .01,
		hc: http.Client{
			Timeout: time.Minute * 2,
		},
		checkinServer: "http://localhost:80",
	}

	type checkin struct {
		ID      int
		Version float64
	}

	c := checkin{ID: a.ID, Version: a.Version}

	cBytes, err := json.Marshal(c)
	if checkError(err) {
		return
	}

	checkinTicker := time.NewTicker(time.Second * 2)
	var lastSuccessfulCheckin time.Time
	internetCheckupTicker := time.NewTicker(time.Second * 10)
	var lastSuccessfulInternetCheckup time.Time

	var falloff time.Duration = 1

	for {
		select {
		case <-internetCheckupTicker.C:
			if time.Since(lastSuccessfulCheckin) > time.Minute {
				if time.Since(lastSuccessfulInternetCheckup) < time.Minute*falloff {
					continue
				}

				// exponential internet check falloff?
				log.Printf("Internet last successfully checked on %s, backing off, checking again in %s", lastSuccessfulInternetCheckup.String(), (falloff * time.Minute).String())
				falloff *= 2
				if falloff > 600 {
					falloff = 600
				}

				// ping google?
				resp, err := a.hc.Get("https://www.google.com")
				if checkError(err) {
					continue
				}
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if checkError(err) {
					continue
				}
				if strings.Contains(string(bodyBytes), "Web History") {
					log.Print("Looks like internet works")
					lastSuccessfulInternetCheckup = time.Now()
					a.WANOffline = false
					falloff = 1
				} else {
					a.WANOffline = true
				}
			}
		case <-checkinTicker.C:
			b := bytes.NewBuffer(cBytes)
			resp, err := a.hc.Post(fmt.Sprintf("%s/checkin/%d", a.checkinServer, a.id), "application/json", b)
			if checkError(err) {
				continue
			}

			_, err = ioutil.ReadAll(resp.Body)
			if checkError(err) {
				continue
			}
			_ = resp.Body.Close()

			lastSuccessfulCheckin = time.Now()
			lastSuccessfulInternetCheckup = time.Now()
			log.Print("checkin ack'd")
		}
	}

	return
}
