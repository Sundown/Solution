package prism

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func IsConstant(e Expression) bool {
	switch e.(type) {
	case Int, Real, Char, Bool:
		return true
	}

	return false
}

func Init(env *Environment) *Environment {
	if len(os.Args) == 1 {
		Error("No files input").Exit()
	}

	var err error
	for i, s := range os.Args {
		if s[0:2] == "--" {
			switch s[2:] {
			case "emit":
				i++
				if len(os.Args) > i {
					switch os.Args[i] {
					case "llvm", "asm", "purellvm":
						env.EmitFormat = os.Args[i]
					default:
						Error("emit expected one of llvm, asm.").Exit()
					}

					Verbose("Emitting", env.EmitFormat)
				} else {
					Error("emit requires argument").Exit()
				}

			case "o":
				if len(os.Args) > i+1 {
					i++
					env.Output = os.Args[i]
				} else {
					Error("output expected filename").Exit()
				}
			case "verbose":
				quietP = false
			case "O":
				i++
				if len(os.Args) > i {
					switch os.Args[i] {
					case "0", "1", "2", "3":
						l, err := strconv.ParseInt(os.Args[i], 10, 32)
						if err != nil {
							Error("optimisation expected integer (0, 1, 3) or \"fast\"").Exit()
						}

						env.Optimisation = &l

					case "fast":
						l := int64(4)
						env.Optimisation = &l
					}
				} else {
					Error("optimisation expected level").Exit()
				}
			}

		} else {
			env.File, err = filepath.Abs(os.Args[1])
			if err != nil {
				Error("Trying to use " + os.Args[1] + " as input file, not found.").Exit()
			}

		}
	}

	return env
}

func Emit(env *Environment) {
	out := []byte((*env.Module).String())

	var sum [32]byte = sha256.Sum256(out)
	temp_name := env.Output + "_" + hex.EncodeToString(sum[:]) + ".ll"

	if env.EmitFormat == "purellvm" {
		ioutil.WriteFile(env.Output+".ll", out, 0644)
		Notify("compiled", env.Output, "to LLVM").Exit()
	} else {
		Verbose("Temp file", temp_name)
		ioutil.WriteFile(temp_name, out, 0644)
	}

	VerifyClangVersion()

	opt := "-Ofast" // TODO change to fast once trap bug fixed
	if env.Optimisation != nil {
		f := strconv.FormatInt(*env.Optimisation, 10)
		Verbose("Optimisation level", f)
		opt = "-O" + f
	}

	sp := ""
	lp := ""
	str := "executable"
	ext := ""

	if env.EmitFormat == "asm" {
		err := exec.Command("clang", temp_name, "-o", env.Output+".s", "-S").Run()
		exec.Command("rm", "-f", temp_name).Run()
		if err != nil {
			Error(err.Error()).Exit()
		}

		Notify("compiled", env.Output, "to Assembly").Exit()
	}

	if env.EmitFormat == "llvm" {
		sp = "-S"
		lp = "-emit-llvm"
		str = "LLVM"
		ext = ".ll"
	} else if env.EmitFormat == "asm" {
		sp = "-S"
		str = "Assembly"
		ext = ".s"
	}

	err := exec.Command("clang", temp_name, opt, sp, lp, "-o", env.Output+ext).Run()
	exec.Command("rm", "-f", temp_name).Run()
	if err != nil {
		Error(err.Error()).Exit()
	} else {
		Notify("compiled", env.Output, "to", str)
	}
}

func PilotEmit(env *Environment) (string, bool) {
	out := []byte((*env.Module).String())
	sum := [32]byte(sha256.Sum256(out))
	temp_name := env.Output + "_" + hex.EncodeToString(sum[:]) + ".ll"

	ioutil.WriteFile(temp_name, out, 0644)

	VerifyClangVersion()

	err := exec.Command("clang", temp_name, "-Og", "-o", env.Output).Run()
	if err != nil {
		return err.Error(), false
	}

	res, err := exec.Command("./" + env.Output).Output()
	exec.Command("rm", "-f", temp_name, env.Output).Run()
	if err != nil {
		return err.Error(), false
	} else {
		return string(res), true
	}

}

type Runtime struct {
	EmitFormat   string
	Output       string
	Verbose      *bool
	Optimisation *int64
	File         string
}

// VerifyClangVersion ensures Clang is installed and at least v12
func VerifyClangVersion() {
	s, err := exec.Command("clang", "--version").Output()

	if err != nil {
		Error("unable to find Clang").Exit()
	}

	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	if !re.MatchString(string(s)) {
		Error("Cannot determine Clang version").Exit()
	}

	ver, err := strconv.ParseFloat(re.FindAllString(string(s), 1)[0], 32)

	if err != nil {
		Panic(err.Error())
	}

	if ver < 12 {
		Error(`Requires clang version 12+`).Exit()
	}

	Verbose("Clang version " + strconv.FormatFloat(ver, 'f', -1, 32))
}

var reset = "\033[0m"

var quietP = true

// Red format for ...strings then reset colour
func Red(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[31m", strings.Join(s, " "), reset)
}

// Green format for ...strings then reset colour
func Green(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[32m", strings.Join(s, " "), reset)
}

// Blue format for ...strings then reset colour
func Blue(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[34m", strings.Join(s, " "), reset)
}

// Yellow format for ...strings then reset colour
func Yellow(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[33m", strings.Join(s, " "), reset)
}

// Bold format for ...strings then reset
func Bold(s ...string) string {
	return fmt.Sprintf("%s%s%s", "\033[1m", strings.Join(s, " "), reset)
}

type er struct {
	code int
}

// Ref gets reference
func Ref(s string) *string {
	// Heaps of them!
	return &s
}

// Verbose message to be displayed depending on quietP
func Verbose(parts ...string) er {
	if !quietP {
		fmt.Println(Blue(" - "), strings.Join(parts, " "))
	}

	return er{code: 0}
}

// Notify general message
func Notify(parts ...string) er {
	fmt.Println(Green(" + "), strings.Join(parts, " "))
	return er{code: 0}
}

// Error message
func Error(parts ...string) er {
	fmt.Println(Red(" ! "), strings.Join(parts, " "))
	return er{code: 1}
}

// Warn message
func Warn(parts ...string) er {
	fmt.Println(Yellow(" ? "), strings.Join(parts, " "))
	return er{code: 1}
}

// Panic prints red and exits
func Panic(err string, args ...interface{}) {
	_, path, line, _ := runtime.Caller(1)
	fmt.Printf(Red(" !  ")+err+"\n", args...)
	fmt.Printf(path+":%d\n", line)
	os.Exit(1)
}

// Exit based on code
func (e er) Exit() {
	os.Exit(e.code)
}
