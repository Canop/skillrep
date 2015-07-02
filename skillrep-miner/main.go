package main

import (
	"log"
	"skillrep/core"
)

func main() {
	ta := core.AnalyzeDatabase()
	log.Printf("Analysis done. %d tokens found", len(ta.Tokens))
	arr := ta.Top(100)
	for i, a := range arr {
		log.Printf("#%d : %s -> %d questions, %d closed, avgBAS=%.2f\n", i,
			a.Token, a.NbQuestions, a.NbClosures,
			float32(a.SumBestAnswerScore)/float32(a.NbQuestions))
	}
}
