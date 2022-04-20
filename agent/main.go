package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type agent struct {
	id            int
	hc            http.Client
	checkinServer string
}

func checkError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

func main() {

	var a *agent

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	a = &agent{
		id: 1,
		hc: http.Client{
			Timeout: time.Minute * 2,
		},
		checkinServer: "http://localhost:80",
	}

	checkinTicker := time.NewTicker(time.Second * 2)

	type checkin struct {
		ID int
	}

	c := checkin{ID: a.id}

	cBytes, err := json.Marshal(c)
	if checkError(err) {
		return
	}

	for {
		select {
		case <-checkinTicker.C:
			fmt.Println("checking in")
			b := bytes.NewBuffer(cBytes)
			resp, err := a.hc.Post(fmt.Sprintf("%s/checkin/%d", a.checkinServer, a.id), "application/json", b)
			if checkError(err) {
				return
			}

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if checkError(err) {
				return
			}
			resp.Body.Close()

			log.Println(string(bodyBytes))
		}
	}

}
