package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	"regexp"
	"skillrep/core"
	"strconv"
	"time"
)

// url with path /dbstats
func handleJsonDBStatsQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	srq := &core.DBStatsQuery{}
	srr := srq.Answer()
	m, _ := json.Marshal(srr)
	w.Write(m)
}

// url with paths:
//           /users
//           /users/123456
func handleJsonUsersQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=120") // 2 minutes
	srq := &core.UsersQuery{}
	srq.Page, _ = strconv.Atoi(r.FormValue("page"))
	srq.Search = r.FormValue("search")
	srq.Fix()
	srr := srq.Answer()
	m, _ := json.Marshal(srr)
	w.Write(m)
}

var userRequestRegex = regexp.MustCompile(`/(\d+)$`)

func handleJsonUserQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Println("User Query - Path:", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=120") // 2 minutes
	w.Header().Set("Access-Control-Allow-Origin", "*")
	m := userRequestRegex.FindStringSubmatch(r.URL.Path)
	if len(m) < 1 {
		http.Error(w, "no user id found in path", 400)
		return
	}
	uq := &core.UserQuery{}
	uq.UserId, _ = strconv.Atoi(m[1])
	ur := uq.Answer()
	b, _ := json.Marshal(ur)
	w.Write(b)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.HandleFunc("/users", handleJsonUsersQuery)
	mux.HandleFunc("/users/", handleJsonUserQuery)
	mux.HandleFunc("/dbstats", handleJsonDBStatsQuery)
	graceful.Run(fmt.Sprintf(":%d", core.Config().Port), 10*time.Second, mux)
}
