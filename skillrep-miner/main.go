package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"skillrep/core"
	"sort"
	"time"
)

func main() {
	start := time.Now()
	ta := core.AnalyzeDatabase()
	fmt.Println("Analysis took ", time.Since(start))
	fmt.Printf("%d questions\n", ta.NbQuestions)
	fmt.Printf("%d bad questions (%.1f%%)\n",
		ta.NbBadQuestions, 100*float32(ta.NbBadQuestions)/float32(ta.NbQuestions))
	fmt.Printf("%d tokens found\n", len(ta.Tokens))
	arr := ta.ToArr(400)
	fmt.Println("Nb tokens in at least 400 questions:", len(arr))
	fmt.Println("Most frequent words:")
	sort.Sort(core.ByNbQuestions(arr))
	mostFrequent := arr[:99]
	sort.Sort(core.ByAverageBestAnswerScore(mostFrequent))
	mostFrequent.PrintMarkdownTable(100)
	fmt.Println("Good Omen:")
	sort.Sort(core.ByAverageBestAnswerScore(arr))
	arr.PrintMarkdownTable(30)
	fmt.Println("Bad Omen:")
	sort.Sort(sort.Reverse(core.ByAverageBestAnswerScore(arr)))
	arr.PrintMarkdownTable(30)

	be := core.NewBayesianEstimator(ta)
	b, _ := json.Marshal(be.Probs)
	err := ioutil.WriteFile("bayes.json", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bayesian Estimator Data written to bayes.json")

	s := "A perfectly ordinary title"
	fmt.Println(s, "=>", be.Estimate(s))
	s = "emacs loop in a monad"
	fmt.Println(s, "=>", be.Estimate(s))
	s = "birt wall nokia monad  combined..."
	fmt.Println(s, "=>", be.Estimate(s))
	s = "wordpress what crystal license"
	fmt.Println(s, "=>", be.Estimate(s))
}
