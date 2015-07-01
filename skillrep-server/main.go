package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	"skillrep/core"
	"strconv"
	"time"
)

func handleJsonDBStatsQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	srq := &core.DBStatsQuery{}
	srr := srq.Answer()
	m, _ := json.Marshal(srr)
	w.Write(m)
}

func handleJsonUsersQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	srq := &core.UsersQuery{}
	srq.Page, _ = strconv.Atoi(r.FormValue("page"))
	srq.Search = r.FormValue("search")
	srq.Fix()
	srr := srq.Answer()
	m, _ := json.Marshal(srr)
	w.Write(m)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.HandleFunc("/users", handleJsonUsersQuery)
	mux.HandleFunc("/dbstats", handleJsonDBStatsQuery)
	graceful.Run(fmt.Sprintf(":%d", core.Config().Port), 10*time.Second, mux)
}
