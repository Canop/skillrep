package core

import (
	"fmt"
	"regexp"
	"strings"
)

var tokenRegex = regexp.MustCompile(`\w{4,}`)

type TokenStats struct {
	Token              string
	NbQuestions        int
	NbClosures         int
	NbBadQuestions     int
	SumBestAnswerScore int
}

func (s TokenStats) AvgBestAnswerScore() float32 {
	return float32(s.SumBestAnswerScore) / float32(s.NbQuestions)
}

type TokenAnalysis struct {
	Tokens         map[string]TokenStats
	NbQuestions    int
	NbClosures     int
	NbBadQuestions int
}

type TokenStatsArray []TokenStats

type ByNbQuestions TokenStatsArray

func (arr ByNbQuestions) Less(i, j int) bool {
	return arr[i].NbQuestions > arr[j].NbQuestions
}
func (arr ByNbQuestions) Len() int {
	return len(arr)
}
func (arr ByNbQuestions) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

type ByAverageBestAnswerScore TokenStatsArray

func (arr ByAverageBestAnswerScore) Less(i, j int) bool {
	return arr[i].AvgBestAnswerScore() > arr[j].AvgBestAnswerScore()
}
func (arr ByAverageBestAnswerScore) Len() int {
	return len(arr)
}
func (arr ByAverageBestAnswerScore) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (ta TokenAnalysis) ToArr(nbQuestionsMin int) TokenStatsArray {
	arr := make([]TokenStats, 0, len(ta.Tokens))
	for _, ts := range ta.Tokens {
		if ts.NbQuestions >= nbQuestionsMin {
			arr = append(arr, ts)
		}
	}
	return arr
}

func (arr TokenStatsArray) PrintMarkdownTable(maxSize int) {
	fmt.Println("|  #  |         Token        | Questions | Closed | Bad Q. |Avg BAS |")
	fmt.Println("|:---:|:--------------------:|-----------|--------|--------|--------|")
	for i, s := range arr {
		if i > maxSize {
			break
		}
		fmt.Printf("|%4d |%-22s|%10d |%6.1f%% |%6.1f%% |%7.1f |\n", i+1, s.Token,
			s.NbQuestions,
			100*float32(s.NbClosures)/float32(s.NbQuestions),
			100*float32(s.NbBadQuestions)/float32(s.NbQuestions),
			s.AvgBestAnswerScore())
	}
}

func AnalyzeDatabase() *TokenAnalysis {
	ta := &TokenAnalysis{}
	ta.Tokens = make(map[string]TokenStats)
	sql := `select title, closeddate,
	(select coalesce(max(score+2*accepted::int),0) from answer where question=question.id)
	from question limit 100000000`
	rows, err := db.TimedQuery(sql)
	die(err)
	var title string
	var closedDate int
	var maxAnswerScore int
	for rows.Next() {
		err = rows.Scan(&title, &closedDate, &maxAnswerScore)
		die(err)
		ta.NbQuestions++
		if closedDate != 0 {
			maxAnswerScore = 0
		} else if maxAnswerScore > 10 {
			maxAnswerScore = 10
		}
		if maxAnswerScore < 4 {
			ta.NbBadQuestions++
		}
		title = strings.ToLower(title)
		tokens := tokenRegex.FindAllString(title, -1)
		for _, t := range tokens {
			ts := ta.Tokens[t]
			ts.Token = t
			ts.NbQuestions++
			if closedDate != 0 {
				ts.NbClosures++
			}
			if maxAnswerScore < 4 {
				ts.NbBadQuestions++
			}
			ts.SumBestAnswerScore += maxAnswerScore
			ta.Tokens[t] = ts
		}
	}
	return ta
}
