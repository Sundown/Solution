package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/llir/llvm/ir"
)

func (rt *Runtime) ParseArgs() {
	if len(os.Args) == 1 {
		Error("No files input").Exit()
	}

	var err error
	for i, s := range os.Args {
		if s[0] == '-' {
			switch s[1:] {
			case "emit":
				if len(os.Args) > i+1 {
					if os.Args[i+1] != "llvm" && os.Args[i+1] != "asm" {
						Error("emit expected one of llvm, asm.").Exit()
					}

					i++
					rt.Emit = os.Args[i]
					Verbose("Emitting", rt.Emit)
				} else {
					Error("emit requires argument").Exit()
				}

			case "o":
				if len(os.Args) > i+1 {
					i++
					rt.Output = os.Args[i]
				} else {
					Error("output expected filename").Exit()
				}
			case "verbose":
				Quietp = false
			}

		} else {
			rt.File, err = filepath.Abs(os.Args[1])
			if err != nil {
				Error("Trying to use " + os.Args[1] + " as input file, not found.").Exit()
			}

		}
	}
}

func (rt *Runtime) HandleEmit(mod *ir.Module) {
	out := []byte(mod.String())

	var sum [32]byte = sha256.Sum256(out)
	temp_name := rt.Output + "_" + hex.EncodeToString(sum[:]) + ".ll"
	Verbose("Temp file", temp_name)

	if rt.Emit == "llvm" {
		ioutil.WriteFile(rt.Output+".ll", out, 0644)
		Notify("Compiled", rt.Output, "to LLVM").Exit()
	} else {
		ioutil.WriteFile(temp_name, out, 0644)
	}

	VerifyClangVersion()

	if rt.Emit == "asm" {
		err := exec.Command("clang", temp_name, "-o", rt.Output+".s", "-S").Run()
		exec.Command("rm", "-f", temp_name).Run()
		if err != nil {
			Error(err.Error()).Exit()
		}

		Notify("Compiled", rt.Output, "to Assembly").Exit()
	}

	err := exec.Command("clang", temp_name, "-o", rt.Output).Run()
	exec.Command("rm", "-f", temp_name).Run()
	if err != nil {
		Error(err.Error()).Exit()
	} else {
		Notify("Compiled", rt.Output, "to executable")
	}
}
