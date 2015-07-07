package core

import (
	"log"
	"math"
	"strings"
)

type BayesianEstimator struct {
	Probs map[string]float64
}

func NewBayesianEstimator(ta *TokenAnalysis) *BayesianEstimator {
	// prB : overall probability that any question is Bad
	prB := float64(ta.NbBadQuestions) / float64(ta.NbQuestions)
	log.Println("prB:", prB)
	arr := ta.ToArr(300) // keeping only the words we find in at least 300 questions
	log.Println("Kept words:", len(arr))
	be := BayesianEstimator{}
	be.Probs = make(map[string]float64)
	for _, ts := range arr {
		ta.Tokens[ts.Token] = ts
		// prTgG : probability to find the token in a good question
		prTgG := float64(ts.NbQuestions-ts.NbBadQuestions) / float64(ta.NbQuestions-ta.NbBadQuestions)
		// prTgB : probability to find the token in a bad question
		prTgB := float64(ts.NbBadQuestions) / float64(ta.NbBadQuestions)
		// prBgT : probability the question is bad given it contains the token
		prBgT := prTgB * prB / (prTgB*prB + prTgG*(1-prB))
		be.Probs[ts.Token] = prBgT
	}
	return &be
}

// Estimates the probability the passed title is of a bad question
func (be BayesianEstimator) Estimate(title string) float64 {
	title = strings.ToLower(title)
	tokens := tokenRegex.FindAllString(title, -1)
	var η float64
	for _, token := range tokens {
		if prBgT, ok := be.Probs[token]; ok {
			log.Println("  ", token, "->", prBgT)
			η += math.Log(1-prBgT) - math.Log(prBgT)
		}
	}
	p := 1 / (1 + math.Exp(η))
	return p
}
