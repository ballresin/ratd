package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type agent struct {
	ID          int
	LastCheckin time.Time
}

func (s *server) getAgentForID(id int) (agent, error) {
	var a agent
	q := "SELECT id, latest_checkin_ts FROM agents WHERE id = ?"
	err := s.db.QueryRowContext(context.Background(), q, id).Scan(&a.ID, &a.LastCheckin)
	if checkError(err) {
		return a, fmt.Errorf("no agent found for id %d %w", id, err)
	}

	return a, nil
}

func (s *server) handleAgentReq(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// handle request for agent data
	id, err := strconv.Atoi(params.ByName("agentID"))
	if checkError(err) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a, err := s.getAgentForID(id)
	if checkError(err) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	aJSONBytes, err := json.Marshal(a)
	if checkError(err) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(aJSONBytes)
}
