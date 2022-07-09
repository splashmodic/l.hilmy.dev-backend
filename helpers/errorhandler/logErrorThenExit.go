package errorhandler

import "fmt"

func LogErrorThenPanic(err *error) {
	panic(fmt.Sprintf("Panic: %v\n", *err))
}
