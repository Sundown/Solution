package compiler

import (
	"sundown/solution/parse"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineReduce(app *parse.Application) value.Value {
	if app.Argument.TypeOf.Tuple == nil {
		panic("Map requires Tuple")
	}

	// Function to be mapped
	fn := app.Argument.Atom.Tuple[0].Atom.Function

	// The vector in AST
	vec := app.Argument.Atom.Tuple[1]

	// The vector in LLVM
	llvec := state.CompileExpression(vec)

	head_type := vec.TypeOf.AsLLType()
	elm_type := vec.Atom.Vector[0].TypeOf.AsLLType()

	if fn.Gives.Equals(parse.AtomicType("Void")) {
		panic("Cannot reduce with void function")
	}

	// Map is 1:1 so leng and cap are just copied from input vector
	leng := state.Block.NewGetElementPtr(head_type, llvec, I32(0), vectorLenOffset)

	accum := state.Block.NewAlloca(elm_type)
	state.Block.NewStore(state.DefaultValue(vec.Atom.Vector[0].TypeOf), accum)

	// -------------
	// # LOOP BODY #
	// -------------
	if app.Argument.Atom.Tuple[1].TypeOf.Vector != nil {
		vec_body := state.Block.NewLoad(
			types.NewPointer(elm_type),
			state.Block.NewGetElementPtr(head_type, llvec, I32(0), vectorBodyOffset))

		counter := state.Block.NewAlloca(types.I64)
		state.Block.NewStore(I64(0), counter)

		// Body
		// Get elem, add to accum, increment counter, conditional jump to body
		loopblock := state.CurrentFunction.NewBlock("")
		state.Block.NewBr(loopblock)
		// Add to accum
		cur_counter := loopblock.NewLoad(types.I64, counter)

		var cur_elm value.Value = loopblock.NewGetElementPtr(elm_type, vec_body, cur_counter)

		if vec.Atom.Vector[0].TypeOf.Atomic != nil {
			cur_elm = loopblock.NewLoad(elm_type, cur_elm)
		}

		call := loopblock.NewCall(
			state.CompileExpression(app.Argument.Atom.Tuple[0]), // fn
			state.CompileDiscreteVector(vec.Atom.Vector[0].TypeOf,
				[]value.Value{
					// These should be replaced with IR calls to GEP and stuff...
					cur_elm,
					state.Block.NewLoad(elm_type, accum),
				}))

		loopblock.NewStore(
			call,
			accum)

		// Increment counter
		loopblock.NewStore(
			loopblock.NewAdd(
				loopblock.NewLoad(types.I64, counter),
				I64(1)),
			counter)

		// Possibly change load to another add or something, probably expensive
		cond := loopblock.NewICmp(enum.IPredSLT,
			loopblock.NewLoad(types.I64, counter),
			loopblock.NewLoad(types.I64, leng))

		exitblock := state.CurrentFunction.NewBlock("")
		loopblock.NewCondBr(cond, loopblock, exitblock)
		state.Block = exitblock

		return state.Block.NewLoad(elm_type, accum)
	} else {
		panic("Map needs (F, [T])")
	}
}

/*
func (state *State) CompileInlineReduce(app *parse.Application) value.Value {
	fn := app.Argument.Atom.Tuple[0]
	llfn := state.CompileExpression(fn)

	arg := app.Argument.Atom.Tuple[1]
	if arg.TypeOf.Vector == nil {
		fmt.Println(arg.TypeOf.String())
		panic("Sum requires Vector")
	}

	typ := arg.TypeOf.Vector

	lltyp := typ.AsLLType()

	vec := arg

	llvec := state.CompileExpression(vec)

	counter := state.Block.NewAlloca(types.I64)
	state.Block.NewStore(I64(0), counter)

	accum := state.Block.NewAlloca(lltyp)
	state.Block.NewStore(state.DefaultValue(typ), accum)

	// Body
	// Get elem, add to accum, increment counter, conditional jump to body
	loopblock := state.CurrentFunction.NewBlock("")
	state.Block.NewBr(loopblock)
	state.Block = loopblock

	// Add to accum
	cur_counter := state.Block.NewLoad(types.I64, counter)

	cur_elm := state.Block.NewLoad(lltyp, state.Block.NewGetElementPtr(lltyp, state.Block.NewLoad(types.NewPointer(lltyp), state.Block.NewGetElementPtr(typ.AsVector().AsLLType(), llvec, I32(0), vectorBodyOffset)), cur_counter))

	// Accum <- accum + current element
	state.Block.NewStore(
		state.Block.NewCall(llfn,
			state.CompileDiscreteVector(typ, []value.Value{
				state.Block.NewLoad(lltyp, accum),
				cur_elm,
			}),
		), accum)

	// Increment counter
	state.Block.NewStore(
		state.Block.NewAdd(state.Block.NewLoad(lltyp, counter), I64(1)),
		counter)

	cond := state.Block.NewICmp(
		enum.IPredSLT,
		cur_counter,
		state.Block.NewLoad(
			types.I64,
			state.Block.NewGetElementPtr(
				typ.AsVector().AsLLType(),
				llvec,
				I32(0),
				vectorLenOffset)))

	exitblock := state.CurrentFunction.NewBlock("")
	state.Block.NewCondBr(cond, state.Block, exitblock)
	state.Block = exitblock

	return state.Block.NewLoad(lltyp, accum)
}
*/
