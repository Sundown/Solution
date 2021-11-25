package compiler

import (
	"sundown/solution/parse"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// Need tuple element access before this is useful, ignore for now
func (state *State) CompileInlineFoldl(app *parse.Application) value.Value {
	fn := app.Argument.Atom.Tuple[0]
	llfn := state.CompileExpression(fn)
	typ := app.Argument.TypeOf.Vector

	lltyp := typ.AsLLType()

	vec := app.Argument

	llvec := state.CompileExpression(vec)

	counter := state.Block.NewAlloca(types.I64)
	state.Block.NewStore(I64(0), counter)

	accum := state.Block.NewAlloca(lltyp)
	state.Block.NewStore(state.Number(typ, 1), accum)

	// Body
	// Get elem, add to accum, increment counter, conditional jump to body

	cond_rhs := state.Block.NewLoad(
		types.I64,
		state.Block.NewGetElementPtr(
			typ.AsVector().AsLLType(),
			llvec,
			I32(0),
			vectorLenOffset))

	// This is a bit messy
	loopblock := state.CurrentFunction.NewBlock("")
	state.Block.NewBr(loopblock)
	state.Block = loopblock
	// ---

	// Add to accum
	cur_counter := loopblock.NewLoad(types.I64, counter)

	// Accum <- accum * current element
	ll_tuple := state.Block.NewAlloca(fn.TypeOf.AsLLType())

	left := accum
	right := loopblock.NewGetElementPtr(
		lltyp,
		state.Block.NewLoad(
			types.NewPointer(lltyp),
			state.Block.NewGetElementPtr(
				typ.AsVector().AsLLType(),
				llvec,
				I32(0),
				vectorBodyOffset)),
		cur_counter)

	state.Block.NewStore(left, state.GEP(ll_tuple, I32(0), I32(0)))
	state.Block.NewStore(right, state.GEP(ll_tuple, I32(0), I32(1)))

	loopblock.NewStore(state.Block.NewCall(llfn, ll_tuple), accum)

	cond := loopblock.NewICmp(
		enum.IPredSLT,
		loopblock.NewAdd(cur_counter, I64(1)),
		cond_rhs)

	// Increment counter
	loopblock.NewStore(
		loopblock.NewAdd(loopblock.NewLoad(types.I64, counter), I64(1)),
		counter)

	exitblock := state.CurrentFunction.NewBlock("")
	loopblock.NewCondBr(cond, loopblock, exitblock)
	state.Block = exitblock

	return state.Block.NewLoad(lltyp, accum)
}
