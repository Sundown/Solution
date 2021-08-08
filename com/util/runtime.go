package util

import (
	"os/exec"
	"regexp"
	"strconv"
)

type Runtime struct {
	Emit   *string
	Output *string
}

func VerifyClangVersion() {
	s, err := exec.Command("clang", "--version").Output()

	if err != nil {
		Error("unable to find Clang")
	}

	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	if !re.MatchString(string(s)) {
		Error("Cannot determine Clang version").Exit()
	}

	ver, err := strconv.ParseFloat(re.FindAllString(string(s), 1)[0], 32)

	if err != nil {
		panic(err)
	}

	if ver < 12 {
		Error(`Requires clang version 12+`).Exit()
	}

	Bene("Using Clang version " + strconv.FormatFloat(ver, 'f', -1, 32))
}
