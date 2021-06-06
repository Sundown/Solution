package main

import (
	"io/ioutil"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Var struct {
	Ident   string
	Mutable bool
	Value   value.Value
}

var (
	module    = ir.NewModule()
	main_func = module.NewFunc("main", types.I32)
	entry     = main_func.NewBlock("entry")
)

func print_module() {
	ioutil.WriteFile("code.ll", []byte(module.String()), 0644)
}

func node_is_number(p *Expression) bool {
	return p.Primary != nil && (p.Primary.Nat != nil || p.Primary.Real != nil)
}

func val_is_number(v value.Value) bool {
	return v.Type() == types.I64 || v.Type() == types.I32
}

func gen(expr *Expression) value.Value {
	if a := expr.Application; a != nil {
		switch a.Op {
		case "return":
			entry.NewRet(gen(a.Atoms[0]))
		case "global":
			temp := gen(a.Atoms[1])
			glob := module.NewGlobal(*a.Atoms[0].Primary.Noun, temp.Type())
			entry.NewStore(temp, glob)
			return entry.NewLoad(temp.Type(), glob)
		case "add":
			lhs := gen(a.Atoms[0])
			if !val_is_number(lhs) {
				panic("add: arg[0] must be number")
			}

			rhs := gen(a.Atoms[1])
			if !val_is_number(rhs) {
				panic("add: arg[1] must be number")
			}

			return entry.NewAdd(lhs, rhs)
		case "sub":
			lhs := gen(a.Atoms[0])
			if !val_is_number(lhs) {
				panic("sub: arg[0] must be number")
			}

			rhs := gen(a.Atoms[1])
			if !val_is_number(rhs) {
				panic("sub: arg[1] must be number")
			}

			return entry.NewSub(lhs, rhs)
		case "mul":
			lhs := gen(a.Atoms[0])
			if !val_is_number(lhs) {
				panic("mul: arg[0] must be number")
			}

			rhs := gen(a.Atoms[1])
			if !val_is_number(rhs) {
				panic("mul: arg[1] must be number")
			}

			return entry.NewMul(lhs, rhs)
		case "div":
			lhs := gen(a.Atoms[0])
			if !val_is_number(lhs) {
				panic("div: arg[0] must be number")
			}

			rhs := gen(a.Atoms[1])
			if !val_is_number(rhs) {
				panic("div: arg[1] must be number")
			}

			return entry.NewFDiv(lhs, rhs)
		}
	} else if p := expr.Primary; p != nil {
		switch {
		case p.Real != nil:
			return constant.NewFloat(types.Double, *p.Real)
		case p.Nat != nil:
			return constant.NewInt(types.I32, int64(*p.Nat))
		case p.String != nil:
			return constant.NewCharArrayFromString(*p.String)
		case p.Noun != nil:
			/* todo */
			return nil
		case p.Bool != nil:
			return constant.NewBool(*p.Bool)
		default:
			return constant.NewNull(types.I32Ptr)
		}
	}

	return nil
}
