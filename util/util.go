package util

import (
	"fmt"
	"strings"
)

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorYellow = "\033[33m"

func Error(parts ...string) {
	fmt.Println("sunday: " + colorRed + "error: " + colorReset + strings.Join(parts, " "))
}

func Warn(parts ...string) {
	fmt.Println("sunday: " + colorYellow + "warning: " + colorReset + strings.Join(parts, " "))
}