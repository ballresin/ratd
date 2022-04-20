package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type agent struct {
	id            int
	hc            http.Client
	checkinServer string
	WANOffline    bool
}

func checkError(err error) bool {
	if err != nil {
		log.Print(err)
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

	type checkin struct {
		ID int
	}

	c := checkin{ID: a.id}

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

}
