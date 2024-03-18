package prism

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	lo "github.com/samber/lo"
)

func IsConstant(e Expression) bool {
	switch e.(type) {
	case Int, Real, Char, Bool:
		return true
	case Vector:
		return true // TODO: not always safe
		// this is causing some issues with vectors of strings
		// may be a root issue in matrices and high order vectors
	}

	return false
}

func Init(env *Environment) *Environment {
	if len(os.Args) == 1 {
		Error("No files input").Exit()
	}
	var other []string
	os.Args = append(os.Args[:1], lo.Filter(os.Args, func(s string, i int) bool {
		if s[0] == '-' || (i > 0 && os.Args[i-1][0] == '-') {
			return true
		}

		other = append(other, s)
		return false
	})...)

	emitFormat := flag.String("emit", "exec", "Emit format")
	optimisationLevel := flag.String("opt", "fast", "Optimisation level")
	flag.BoolVar(&env.Verbose, "verbose", false, "Verbose")

	flag.Parse()

	switch *emitFormat {
	case "purellvm", "llvm", "asm", "exec":
		env.EmitFormat = *emitFormat
	}

	switch *optimisationLevel {
	case "0", "1", "2", "3", "fast":
		env.Optimisation = *optimisationLevel
	}

	f, err := filepath.Abs(other[1])
	if err != nil {
		Error("Trying to use " + os.Args[1] + " as input file, not found.").Exit()
	} else {
		env.File = f
	}

	return env
}

// Emit observes emit specifier and writes LLVM IR to file and invokes Clang
// Will exit.
func Emit(env *Environment) {
	out := []byte((*env.Module).String())

	if env.EmitFormat == "purellvm" {
		os.WriteFile(env.Output+".ll", out, 0644)
		Notify("compiled", env.Output, "to LLVM").Exit()
	}

	var sum [32]byte = sha256.Sum256(out) // collision perhaps, error even?
	temp_name := env.Output + "_" + hex.EncodeToString(sum[:]) + ".ll"

	Verbose("Temp file", temp_name)
	os.WriteFile(temp_name, out, 0644)

	VerifyClangVersion()

	Verbose("Optimisation level", env.Optimisation)

	var str string
	var args []string

	switch env.EmitFormat {
	case "asm":
		str = "Assembly"
		args = []string{"-S", "-o", env.Output + ".s"}
	case "llvm":
		str = "LLVM IR"
		args = []string{"-S", "-emit-llvm", "-o", env.Output + ".ll"}
	case "exec":
		str = "Executable"
		args = []string{"-o", env.Output}
	}

	err := exec.Command("clang", append(args, temp_name, "-O"+env.Optimisation)...).Run()

	exec.Command("rm", "-f", temp_name).Run()

	if err != nil {
		Error("Clang returned non-zero exit code")
		Error(err.Error()).Exit()
	}

	Notify("compiled", env.Output, "to", str).Exit()
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
