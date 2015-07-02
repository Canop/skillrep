package core

import (
	"regexp"
	"sort"
	"strings"
)

var tokenRegex = regexp.MustCompile(`\w{4,}`)

type TokenStats struct {
	Token              string
	NbQuestions        int
	NbClosures         int
	SumBestAnswerScore int
}

type TokenAnalysis struct {
	Tokens map[string]TokenStats
}

type TokenStatsArray []TokenStats

func (arr TokenStatsArray) Len() int {
	return len(arr)
}
func (arr TokenStatsArray) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}
func (arr TokenStatsArray) Less(i, j int) bool {
	return arr[i].NbQuestions > arr[j].NbQuestions
}

func (ta TokenAnalysis) ToArr(nbQuestionsMin int) []TokenStats {
	arr := make([]TokenStats, 0, len(ta.Tokens))
	for _, ts := range ta.Tokens {
		if ts.NbQuestions >= nbQuestionsMin {
			arr = append(arr, ts)
		}
	}
	return arr
}
func (ta TokenAnalysis) Top(n int) []TokenStats {
	arr := make([]TokenStats, len(ta.Tokens))
	i := 0
	for _, ts := range ta.Tokens {
		arr[i] = ts
		i++
	}
	sort.Sort(TokenStatsArray(arr))
	return arr[0:n]
}

func AnalyzeDatabase() *TokenAnalysis {
	ta := &TokenAnalysis{}
	ta.Tokens = make(map[string]TokenStats)
	sql := `select title, closeddate,
		(select coalesce(max(score),0) from answer where question=question.id)
		from question limit 1000000`
	rows, err := db.TimedQuery(sql)
	die(err)
	var title string
	var closedDate int
	var maxAnswerScore int
	for rows.Next() {
		err = rows.Scan(&title, &closedDate, &maxAnswerScore)
		die(err)
		title = strings.ToLower(title)
		tokens := tokenRegex.FindAllString(title, -1)
		for _, t := range tokens {
			ts := ta.Tokens[t]
			ts.Token = t
			ts.NbQuestions++
			if maxAnswerScore > 10 {
				maxAnswerScore = 10
			}
			ts.SumBestAnswerScore += maxAnswerScore
			if closedDate != 0 {
				ts.NbClosures++
			}
			ta.Tokens[t] = ts
		}
	}
	return ta
}
