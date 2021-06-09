package codegen

import (
	"fmt"
	"io/ioutil"

	"sundown/girl/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var (
	param  = "x"
	int_bt types.IntType
	int_t  types.Type
	real_t types.Type
	str_t  types.Type
)

type State struct {
	module   *ir.Module
	function *ir.Func
	block    *ir.Block
	fns      map[string]*ir.Func
}

func StartCompiler(path string, block *parser.Program) error {
	state := State{}
	state.module = ir.NewModule()
	state.fns = make(map[string]*ir.Func)

	int_bt = types.IntType{TypeName: "int", BitSize: 64}
	int_t = state.module.NewTypeDef("int", types.I64)
	real_t = state.module.NewTypeDef("real", types.Double)
	str_t = state.module.NewTypeDef("string", types.I8Ptr)

	state.BuiltinDouble()
	state.BuiltinPuts()

	for _, expr := range block.Expression {
		state.compile(expr)
	}

	ioutil.WriteFile(path, []byte(state.module.String()), 0644)

	return nil
}

func (state *State) compile(expr *parser.Expression) value.Value {
	if expr.FnDecl != nil {
		fn := state.module.NewFunc(
			*expr.FnDecl.Ident.Ident,
			GenType(expr.FnDecl.Type.Gives),
			ir.NewParam("x", GenType(expr.FnDecl.Type.Takes)))

		state.function = fn
		state.block = state.function.NewBlock("entry")

		// Step through and codegen each expression in the function until ";"
		for _, expr := range expr.FnDecl.Block.Expression {
			state.compile(expr)
		}

		if state.function.Sig.RetType == types.Void {
			state.block.NewRet(nil)
		}

		state.fns[*expr.FnDecl.Ident.Ident] = fn
		// Constructing this function is over so clear state
		state.block = nil
		state.function = nil
	} else if expr.Primary != nil {
		if expr.Primary.Param != nil {
			return state.function.Params[0]
		} else if expr.Primary.Int != nil {
			return constant.NewInt(&int_bt, *expr.Primary.Int)
		} else if expr.Primary.Real != nil {
			return constant.NewFloat(types.Double, *expr.Primary.Real)
		} else if expr.Primary.Bool != nil {
			if *expr.Primary.Bool == "true" {
				return constant.NewBool(true)
			} else if *expr.Primary.Bool == "false" {
				return constant.NewBool(false)
			}
		} else if expr.Primary.Vec != nil {
			vec, _ := state.compile_vector(expr.Primary.Vec)
			return vec
		} else if expr.Primary.String != nil {
			s := (*expr.Primary.String)[1:len(*expr.Primary.String)-1] + "\x00"
			slen := int64(len(s))
			i := constant.NewCharArrayFromString(s)
			str := state.module.NewGlobalDef("", i)
			ptr := constant.NewGetElementPtr(
				types.NewArray(uint64(slen), types.I8),
				str,
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, 0))

			return ptr
		}
	} else if expr.Application != nil {
		switch *expr.Application.Op.Ident {
		case "return":
			if state.function.Sig.RetType == types.Void {
				state.block.NewRet(nil)
			} else {
				state.block.NewRet(state.compile(expr.Application.Atoms))
			}

		case "head":
			vec, vec_type := state.compile_vector(expr.Application.Atoms.Primary.Vec)
			return state.block.NewLoad(
				types.I32,
				state.block.NewGetElementPtr(
					vec_type,
					vec,
					constant.NewInt(types.I32, 0),
					constant.NewInt(types.I32, 0)))
		default:
			return state.block.NewCall(
				state.fns[*expr.Application.Op.Ident],
				state.compile(expr.Application.Atoms))
		}
	}

	return nil
}

func (state *State) compile_vector(vector []*parser.Expression) (value.Value, *types.VectorType) {
	elm_type := GenPrimaryType(vector[0].Primary)
	fmt.Println(elm_type)
	vec_type := &types.VectorType{
		TypeName: "",
		Scalable: true,
		Len:      uint64(len(vector)),
		ElemType: elm_type}
	vec := state.block.NewAlloca(vec_type)

	for i, elm := range vector {
		ptr := state.block.NewGetElementPtr(
			vec_type,
			vec,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)))
		ptr.InBounds = true

		state.block.NewStore(state.compile(elm), ptr)
	}

	return vec, vec_type
}

func GenType(t *parser.TypeName) types.Type {
	switch *t.Type {
	case "int":
		return int_t
	case "real":
		return real_t
	case "bool":
		return types.I1
	case "void":
		return types.Void
	case "str":
		return str_t
	default:
		return types.Void
	}
}

func GenPrimaryType(p *parser.Primary) types.Type {
	if p != nil {
		switch {
		case p.Int != nil:
			fmt.Println("nice")
			return types.I32
		case p.Real != nil:
			return types.Double
		case p.Bool != nil:
			return types.I1
		case p.String != nil:
			return &types.ArrayType{}
		}
	}

	return nil
}

func gen_literal(p *parser.Primary) value.Value {
	switch {
	case p.Int != nil:
		return constant.NewInt(types.I64, *p.Int)
	case p.Real != nil:
		return constant.NewFloat(types.Double, *p.Real)
	case p.Bool != nil:
		if *p.Bool == "true" {
			return constant.NewBool(true)
		} else if *p.Bool == "false" {
			return constant.NewBool(false)
		}
	case p.String != nil:
		return constant.NewCharArrayFromString(*p.String)
	}

	return nil
}
