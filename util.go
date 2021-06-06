package main

import (
	"fmt"
	"strings"
)

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorYellow = "\033[33m"

func pretty_error(parts ...string) {
	fmt.Println("rib: " + colorRed + "error: " + colorReset + strings.Join(parts, " "))
}

func pretty_warn(parts ...string) {
	fmt.Println("rib: " + colorYellow + "warning: " + colorReset + strings.Join(parts, " "))
}
