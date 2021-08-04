package util

import (
	"fmt"
	"os"
	"strings"
)

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorYellow = "\033[33m"
var bold = "\033[1m"

type E struct {
	code int
}

// Heaps of them!
func Ref(s string) *string {
	return &s
}

func Error(parts ...string) E {
	fmt.Println(bold + "Solution: " + colorRed + "error: " + colorReset + strings.Join(parts, " "))
	return E{code: 1}
}

func (e E) Exit() {
	os.Exit(e.code)
}

func Warn(parts ...string) E {
	fmt.Println(bold + "Solution: " + colorYellow + "warning: " + colorReset + strings.Join(parts, " "))
	return E{code: 1}
}
