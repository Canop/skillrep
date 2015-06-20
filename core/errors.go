package core

import (
	"log"
	"runtime/debug"
)

// die is an acronym for "die if error"
func die(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
}
