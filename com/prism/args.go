package prism

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func Init(env *Environment) *Environment {
	if len(os.Args) == 1 {
		Error("No files input").Exit()
	}

	var err error
	for i, s := range os.Args {
		if s[0:2] == "--" {
			switch s[2:] {
			case "emit":
				if len(os.Args) > i+1 {
					if os.Args[i+1] != "llvm" && os.Args[i+1] != "asm" {
						Error("emit expected one of llvm, asm.").Exit()
					}

					i++
					env.EmitFormat = os.Args[i]
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
				Quietp = false
			case "optimisation", "optimization":
				if len(os.Args) > i+1 {
					i++
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
	Verbose("Temp file", temp_name)

	if env.EmitFormat == "llvm" {
		ioutil.WriteFile(env.Output+".ll", out, 0644)
		Notify("Compiled", env.Output, "to LLVM").Exit()
	} else {
		ioutil.WriteFile(temp_name, out, 0644)
	}

	VerifyClangVersion()

	if env.EmitFormat == "asm" {
		err := exec.Command("clang", temp_name, "-o", env.Output+".s", "-S").Run()
		exec.Command("rm", "-f", temp_name).Run()
		if err != nil {
			Error(err.Error()).Exit()
		}

		Notify("Compiled", env.Output, "to Assembly").Exit()
	}

	opt := ""
	if env.Optimisation != nil {
		f := strconv.FormatInt(*env.Optimisation, 10)
		Verbose("Optimisation level", f)
		opt = "-O" + f
	}

	err := exec.Command("clang", temp_name, opt, "-o", env.Output).Run()
	exec.Command("rm", "-f", temp_name).Run()
	if err != nil {
		Error(err.Error()).Exit()
	} else {
		Notify("Compiled", env.Output, "to executable")
	}
}