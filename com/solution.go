package main

import (
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
	util.VerifyClangVersion()

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
				} else {
					util.Error("emit expected one of [bc, llvm, binary, asm]").Exit()
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
	r, err := os.Open(rt.File)

	if err != nil {
		util.Error(err.Error()).Exit()
	}

	err = lex.Parser.Parse(rt.File, r, l)
	r.Close()

	p := &parse.State{}

	t := p.Parse(l)
	c := &compiler.State{Runtime: rt}

	if rt.Output == "" {
		rt.Output = *p.PackageIdent
	}

	c.Compile(t)

	_, err = exec.Command("clang", rt.Output+".ll", "-o", rt.Output).Output()
	_, err = exec.Command("rm", "-f", rt.Output+".ll").Output()
	if err != nil {
		util.Error(err.Error()).Exit()
	} else {
		util.Notify(rt.Output, "compiled")
	}
}
