package main

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type checkin struct {
	agent

	// use this to block http.ResponseWriter
	wait chan []byte
}

func newCheckin() checkin {
	var c checkin
	c.wait = make(chan []byte, 1)
	return c
}

func (s *server) checkinProcessor() {

	var checkinPool []checkin

	t := time.NewTicker(time.Second * 10)

	for {
		select {
		case a := <-s.checkinChan:
			// handle checkin write to DB
			checkinPool = append(checkinPool, a)
			a.wait <- nil
		case _ = <-t.C:
			// commit checkins to DB
			log.Printf("Committing agent checkins to DB")
			for _, a := range checkinPool {

				q := "UPDATE agents SET latest_checkin_ts = NOW() WHERE id = ?"
				_, err := s.db.ExecContext(context.Background(), q, a.ID)
				if checkError(err) {
					return
				}
			}
		}
	}
}

func (s *server) handleCheckin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// handle agent checkins

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if checkError(err) {
		return
	}

	c := newCheckin()

	err = json.Unmarshal(bodyBytes, &c)
	if checkError(err) {
		return
	}

	log.Printf("Received checkin from agent %d", c.ID)

	r.Body.Close()

	s.checkinChan <- c

	retBytes := <-c.wait

	w.Write(retBytes)

	// we want to optionally write back to agent any pending commands
	// ask cache if commands are queued, deliver if so
	// to do this, we must maintain an open connection until our processor returns a command or nothing

}
