package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"skillrep/core"
)

func handleJsonQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	srq := &core.SRQuery{}
	srq.Fix()
	srr := srq.Answer()
	m, _ := json.Marshal(srr)
	w.Write(m)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/users", handleJsonQuery)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", core.Config().Port), nil))
}
