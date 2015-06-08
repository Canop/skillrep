package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
)

const apiurl = "http://api.stackexchange.com/2.2"

// structure of a StackExchange response
// see https://api.stackexchange.com/docs/wrapper
type Response struct {
	Backoff        int    `json:"backoff"`
	ErrorId        int    `json:"error_id"`
	ErrorMessage   string `json:"error_message"`
	ErrorName      string `json:"error_name"`
	HasMore        bool   `json:"has_more"`
	Page           int    `json:"page"`
	PageSize       int    `json:"page_size"`
	QuotaMax       int    `json:"quota_max"`
	QuotaRemaining int    `json:"quota_remaining"`
	Total          int    `json:"total"`
	Type           string `json:"type"`
}

type Question struct {
	Id int64 `json:"question_id"`
	Title string `json:"title"`
	Tags []string `json:"tags"`
	CreationDate int64 `json:"creation_date"`
	Owner ShallowUser `json:"owner"`
}

type QuestionResponse struct {
	Response
	Items []Question `json:"items"`
}

type Answer struct {
	Id int64 `json:"answer_id"`
	CreationDate int64 `json:"creation_date"`
	IsAccepted bool `json:"is_accepted"`
	Owner ShallowUser `json:"owner"`
	Question Id int `json:"question_id"`
	Score int `json:"score"`
}

type ShallowUser struct {
	Id int64 `json:"user_id"`
	Name string `json:"display_name"`
}

func getQuestion(site string, id int) (*QuestionResponse, error) {
	filter := "!L_Zm1rmoFy)u)LqgLTvHLi"
	rtype := "questions"
	httpClient := new(http.Client)
	url := fmt.Sprintf("%s/%s/%d?site=%s&filter=%s", apiurl, rtype, id, site, filter)
	log.Println("URL: " + url)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r QuestionResponse
	decoder := json.NewDecoder(bufio.NewReader(resp.Body))
	err = decoder.Decode(&r)
	//bytes, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return nil, err
	//}
	//log.Println(string(bytes))
	//err = json.Unmarshal(bytes, &r)
	//if err != nil {
	//	return nil, err
	//}
	return &r, err
}

func main() {
	r, err := getQuestion("stackoverflow", 11353679)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("r:%+v\n", r)
}
