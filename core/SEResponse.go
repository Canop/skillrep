// Stack Exchange Response
package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const apiurl = "http://api.stackexchange.com/2.2"

// structure of a StackExchange response
// see https://api.stackexchange.com/docs/wrapper
type SEResponse struct {
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

type QuestionsSEResponse struct {
	SEResponse
	Questions []*Question `json:"items"`
}

type Question struct {
	Id           int64       `json:"question_id"`
	Title        string      `json:"title"`
	Tags         []string    `json:"tags"`
	CreationDate int64       `json:"creation_date"`
	ClosedDate   int64       `json:"closed_date"`
	Owner        ShallowUser `json:"owner"`
	Answers      []*Answer   `json:"answers"`
}

type Answer struct {
	Id           int64       `json:"answer_id"`
	CreationDate int64       `json:"creation_date"`
	IsAccepted   bool        `json:"is_accepted"`
	Owner        ShallowUser `json:"owner"`
	QuestionId   int         `json:"question_id"`
	Score        int         `json:"score"`
}

type ShallowUser struct {
	Id      int64  `json:"user_id"`
	Name    string `json:"display_name"`
	Profile string `json:"profile_image"`
}

func GetQuestion(site string, id int) (*QuestionsSEResponse, error) {
	filter := "!OfYUQYtgCaZ9JBeJyrvLd85AXer(WSNHQacu))0iZzl"
	httpClient := new(http.Client)
	url := fmt.Sprintf("%s/%s/%d?site=%s&filter=%s", apiurl, "questions", id, site, filter)
	log.Println("URL: " + url)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r QuestionsSEResponse
	decoder := json.NewDecoder(bufio.NewReader(resp.Body))
	err = decoder.Decode(&r)
	return &r, err
}

func GetQuestions(site string, fromDate, toDate int64, page int) (*QuestionsSEResponse, error) {
	//filter := "!OfYUQYtgCaZ9JBeJyrvLd85AXer(WSNHQacu))0iZzl"
	filter := "!OfYUNMJh9fuDMgvKdY4azSGUzEClDUc-_K.I(mDC4F3"
	httpClient := new(http.Client)
	url := fmt.Sprintf("%s/%s?site=%s&filter=%s&sort=creation&order=asc&key=%s", apiurl, "questions", site, filter, config.ApiKey)
	if page > 0 {
		url += fmt.Sprintf("&page=%d", page)
	}
	if fromDate > 0 {
		url += fmt.Sprintf("&fromdate=%d", fromDate)
	}
	if toDate > 0 {
		url += fmt.Sprintf("&todate=%d", toDate)
	}
	log.Println("URL: " + url)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r QuestionsSEResponse
	decoder := json.NewDecoder(bufio.NewReader(resp.Body))
	err = decoder.Decode(&r)
	if r.ErrorMessage != "" {
		log.Printf("ErrorMessage: %s\n", r.ErrorMessage)
	}
	if r.ErrorName != "" {
		log.Fatal("ErrorName:", r.ErrorName)
	}
	return &r, err
}
