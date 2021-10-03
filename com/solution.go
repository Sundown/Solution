package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sundown/solution/compiler"
	"sundown/solution/lex"
	"sundown/solution/parse"
	"sundown/solution/util"
)

func main() {
	rt := &util.Runtime{}
	util.Notify("Solution init...")

	if len(os.Args) == 1 {
		util.Error("No files input").Exit()
	}

	var err error
	for i, s := range os.Args {
		if s[0] == '-' {
			switch s[1:] {
			case "emit":
				if len(os.Args) > i+1 {
					i++
					rt.Emit = os.Args[i]
					util.Verbose("Emitting", rt.Emit)
				} else {
					util.Error("emit expected one of llvm, asm.").Exit()
				}

			case "o":
				if len(os.Args) > i+1 {
					i++
					rt.Output = os.Args[i]
				} else {
					util.Error("output expected filename").Exit()
				}
			}
		} else {
			rt.File, err = filepath.Abs(os.Args[1])
			if err != nil {
				util.Error("Trying to use " + os.Args[1] + " as input file, not found.").Exit()
			}

		}
	}

	l := &lex.State{}
	util.Verbose("Opening", rt.File)
	r, err := os.Open(rt.File)
	defer r.Close()

	if err != nil {
		util.Error(err.Error()).Exit()
	}

	util.Verbose("Init lexer")
	err = lex.Parser.Parse(rt.File, r, l)
	r.Close()

	p := &parse.State{}

	util.Verbose("Init parser")
	t := p.Parse(l)
	c := &compiler.State{Runtime: rt}

	if rt.Output == "" {
		rt.Output = *p.PackageIdent
	}

	util.Verbose("Init compiler")
	mod := c.Compile(t)

	out := []byte(mod.String())

	var sum [32]byte = sha256.Sum256(out)
	temp_name := rt.Output + "_" + hex.EncodeToString(sum[:]) + ".ll"
	util.Verbose("Temp file", temp_name)

	if rt.Emit == "llvm" {
		ioutil.WriteFile(rt.Output+".ll", out, 0644)
		util.Notify("Compiled", rt.Output, "to LLVM").Exit()
	} else {
		ioutil.WriteFile(temp_name, out, 0644)
	}

	util.VerifyClangVersion()

	if rt.Emit == "asm" {
		err = exec.Command("clang", temp_name, "-o", rt.Output+".s", "-S").Run()
		exec.Command("rm", "-f", temp_name).Run()
		if err != nil {
			util.Error(err.Error()).Exit()
		}

		util.Notify("Compiled", rt.Output, "to Assembly").Exit()
	}

	err = exec.Command("clang", "-03", temp_name, "-o", rt.Output).Run()
	exec.Command("rm", "-f", temp_name).Run()
	if err != nil {
		util.Error(err.Error()).Exit()
	} else {
		util.Notify("Compiled", rt.Output, "to executable")
	}
}
