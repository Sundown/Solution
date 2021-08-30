package util

import (
	"fmt"
	"os"
	"strings"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Bold = "\033[1m"

type E struct {
	code int
}

func Ref(s string) *string {
	// Heaps of them!
	return &s
}

func Bene(parts ...string) {
	fmt.Printf("%s[ ]%s %s\n", Green, Reset, strings.Join(parts, " "))
}

func Error(parts ...string) E {
	fmt.Printf("%s[X]%s %s\n", Red, Reset, strings.Join(parts, " "))
	return E{code: 1}
}

func Warn(parts ...string) E {
	fmt.Printf("%s[-]%s %s\n", Yellow, Reset, strings.Join(parts, " "))
	return E{code: 1}
}

func (e E) Exit() {
	os.Exit(e.code)
}
