package errorhandler

import (
	"log"
	"runtime/debug"
)

func LogErrorThenContinue(err *error) {
	log.Printf("Error: %v\n\n\n", *err)
	debug.PrintStack()
}
