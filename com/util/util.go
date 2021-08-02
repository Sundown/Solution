package util

import (
	"fmt"
	"strings"
)

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorYellow = "\033[33m"

// Heaps of them!
func Ref(s string) *string {
	return &s
}

func Error(parts ...string) {
	fmt.Println("solution: " + colorRed + "error: " + colorReset + strings.Join(parts, " "))
}

func Warn(parts ...string) {
	fmt.Println("solution: " + colorYellow + "warning: " + colorReset + strings.Join(parts, " "))
}
