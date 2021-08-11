package util

import (
	"fmt"
	"os"
	"strings"
)

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var bold = "\033[1m"

type E struct {
	code int
}

var counter = 0

// Heaps of them!
func Ref(s string) *string {
	return &s
}

func Bene(parts ...string) {
	counter++
	fmt.Printf("%04d   %s[OKAY]%s   %s\n", counter, colorGreen, colorReset, strings.Join(parts, " "))
}

func Error(parts ...string) E {
	counter++
	fmt.Printf("%04d   %s[ERRO]%s   %s\n", counter, colorRed, colorReset, strings.Join(parts, " "))
	return E{code: 1}
}

func (e E) Exit() {
	os.Exit(e.code)
}

func Warn(parts ...string) E {
	fmt.Printf("%04d   %s[WARN]%s   %s\n", counter, colorYellow, colorReset, strings.Join(parts, " "))
	return E{code: 1}
}
