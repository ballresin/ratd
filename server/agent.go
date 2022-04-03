package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type agent struct {
	ID int
}

func handleAgentReq(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// handle agent checkins

}
