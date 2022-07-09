package random

import (
	"math/rand"
	"strings"
	"time"
)

var charString = "abcdefghijklmnopqrstuvwxyz"
var numberString = "0123456789"

type Options struct {
	WithLowerChar bool
	WithUpperChar bool
	WithNumber    bool
}

func String(n int, options *Options) string {
	chars := ""
	if options.WithLowerChar {
		chars += charString
	}
	if options.WithUpperChar {
		chars += strings.ToUpper(charString)
	}
	if options.WithNumber {
		chars += numberString
	}

	rand.Seed(time.Now().UnixNano())

	randString := ""

	for i := 0; i < n; i++ {
		randIndex := rand.Intn(len(chars))
		randString += string(chars[randIndex])
	}

	return randString
}
