package codegen

import (
	"fmt"
	"io/ioutil"
	"os"
	"sundown/sunday/parser"
	"sundown/sunday/util"

	"github.com/enescakir/emoji"
	llir "github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var (
	param  = "Param"
	int_bt types.IntType
	int_t,
	real_t,
	//types.I8Ptr,
	bool_t types.Type
)

type State struct {
	module   *llir.Module
	function *llir.Func
	Entry    *llir.Func
	block    *llir.Block
	fns      map[string]*llir.Func
}

func StartCompiler(path string, block *parser.Program) error {
	state := State{}
	state.module = llir.NewModule()
	state.fns = make(map[string]*llir.Func)

	int_bt = types.IntType{TypeName: "", BitSize: 64}
	int_t = types.I64
	real_t = state.module.NewTypeDef("Real", types.Double)
	//types.I8Ptr = state.module.NewTypeDef("String", types.I8Ptr)
	bool_t = state.module.NewTypeDef("Bool", types.I1)

	state.BuiltinPuts()
	state.BuiltinDouble()
	state.BuiltinCalloc()

	for _, he := range block.Statements {
		if he.FnDecl != nil {
			state.CompileFunction(he.FnDecl)
		}
	}

	for _, he := range block.Statements {
		if he.Directive != nil {
			state.Direct(he.Directive)
		}
	}

	// Generate entry point
	state.function = state.module.NewFunc("main", types.I32)
	state.block = state.function.NewBlock("entry")
	state.block.NewCall(state.Entry)
	state.block.NewRet(constant.NewInt(types.I32, 0))

	if path == "" {
		if len(packagename) != 0 {
			path = packagename
		} else {
			path = state.Entry.Name()
		}
	}

	ioutil.WriteFile(path+".ll", []byte(state.module.String()), 0644)

	fmt.Println(string(emoji.Dove), " Compiled", path, "successfully")

	return nil
}

func (state *State) CompileFunction(fn *parser.FnDecl) {
	if takes := MakeType(fn.Takes); takes != types.Void {
		state.function = state.module.NewFunc(
			*fn.Ident,
			MakeType(fn.Gives),
			llir.NewParam(param, takes))
	} else {
		state.function = state.module.NewFunc(
			*fn.Ident,
			MakeType(fn.Takes))
	}

	state.block = state.function.NewBlock("entry")

	// Step through and codegen each expression in the function until ";"
	for _, expr := range fn.Expressions {
		state.Compile(expr)
	}

	if state.function.Sig.RetType == types.Void {
		state.block.NewRet(nil)
	}

	state.fns[*fn.Ident] = state.function
	// Constructing this function is over so clear state
	state.block = nil
	state.function = nil
}

func (state *State) Compile(expr *parser.Expression) (value.Value, types.Type) {
	if expr.Primary != nil {
		return state.MakePrimary(expr.Primary)
	} else if expr.Application != nil {
		switch *expr.Application.Function {
		case "Return":
			if state.function.Sig.RetType == types.Void {
				state.block.NewRet(nil)
			} else {
				v, _ := state.Compile(expr.Application.Parameter)
				state.block.NewRet(v)
			}

		case "Head":
			vec, typ := state.Compile(expr.Application.Parameter)

			return state.block.NewLoad(typ, state.block.NewLoad(
				types.NewPointer(typ),
				state.block.NewGetElementPtr(
					BuildVectorType(typ),
					vec,
					constant.NewInt(types.I32, 0),
					constant.NewInt(types.I32, 2)))), BuildVectorType(typ)
		default:
			fn, err := state.fns[*expr.Application.Function]
			if !err {
				util.Error("Function not found")
				os.Exit(1)
			}

			v, t := state.Compile(expr.Application.Parameter)

			return state.block.NewCall(fn, v), t /* TODO: probably wrong */
		}
	}

	return nil, nil
}

func MakeType(t *parser.Type) types.Type {
	switch {
	case t.Primative != nil:
		return NameToType(t.Primative)
	case t.Vector != nil:
		return BuildVectorType(MakeType(t.Vector))
	case t.Struct != nil:
		panic("Struct types not implemented yet")
	default:
		panic("Unknown type class")
	}
}

func NameToType(t *parser.TypeName) types.Type {
	switch *t.Type {
	case "Int":
		return int_t
	case "Real":
		return real_t
	case "Bool":
		return types.I1
	case "Void":
		return types.Void
	case "Str":
		return types.I8Ptr
	default:
		return types.Void
	}
}

func GenPrimaryType(p *parser.Primary) types.Type {
	if p != nil {
		switch {
		case p.Vec != nil:
			if p.Vec[0].Primary == nil {
				panic("Cannot yet calculate type of complex expression inside sub-vector or something")
			}
			/* TODO: this works for some things and breaks others, shouldn't be pointer or something */
			return types.NewPointer(BuildVectorType(GenPrimaryType(p.Vec[0].Primary)))
		case p.Int != nil:
			return int_t
		case p.Real != nil:
			return real_t
		case p.Bool != nil:
			return bool_t
		case p.String != nil:
			return types.I8Ptr
		}
	}

	return nil
}
