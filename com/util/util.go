package util

import (
	"fmt"
	"os"
	"strings"
)

var Reset = "\033[0m"

var Quietp = true

func Red(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[31m", strings.Join(s, " "), Reset)
}

func Green(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[32m", strings.Join(s, " "), Reset)
}

func Blue(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[34m", strings.Join(s, " "), Reset)
}

func Yellow(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[33m", strings.Join(s, " "), Reset)
}

func Bold(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[1m", strings.Join(s, " "), Reset)
}

type E struct {
	code int
}

func Ref(s string) *string {
	// Heaps of them!
	return &s
}

func Verbose(parts ...string) E {
	if !Quietp {
		fmt.Println(Blue(" - "), strings.Join(parts, " "))
	}

	return E{code: 0}
}

func Notify(parts ...string) E {
	fmt.Println(Green(" + "), strings.Join(parts, " "))
	return E{code: 0}
}

func Error(parts ...string) E {
	fmt.Println(Red(" ! "), strings.Join(parts, " "))
	return E{code: 1}
}

func Warn(parts ...string) E {
	fmt.Println(Yellow(" ? "), strings.Join(parts, " "))
	return E{code: 1}
}

func (e E) Exit() {
	os.Exit(e.code)
}
