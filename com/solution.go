package main

import (
	"os"
	"sundown/solution/compiler"
	"sundown/solution/lex"
	"sundown/solution/parse"
	"sundown/solution/util"
)

func main() {
	util.Bene("Solution version " + util.GetSolutionVersion())
	util.VerifyClangVersion()

	rt := &util.Runtime{}
	if len(os.Args) == 1 {
		util.Error("No files input").Exit()
	}

	for i, s := range os.Args {
		if s[0] == '-' {
			switch s[1:] {
			case "emit":
				if len(os.Args) > i+1 {
					i++
					rt.Emit = util.Ref(os.Args[i])
				} else {
					util.Error("emit expected one of [bc, llvm, binary, asm]").Exit()
				}

			case "o":
				if len(os.Args) > i+1 {
					i++
					rt.Output = util.Ref(os.Args[i])
				} else {
					util.Error("output expected filename").Exit()
				}
			}
		}
	}

	l := &lex.State{}
	p := &parse.State{}
	c := &compiler.State{}

	c.Compile(p.Parse(l.Lex(os.Args[1])))
}
