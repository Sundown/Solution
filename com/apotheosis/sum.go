package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileInlineSum(val Value) value.Value {
	typ := val.Type.(prism.VectorType).Type
	lltyp := typ.Realise()

	counter := env.Block.NewAlloca(types.I64)
	env.Block.NewStore(I64(0), counter)

	accum := env.Block.NewAlloca(lltyp)
	env.Block.NewStore(env.DefaultValue(typ), accum)

	// Body
	// Get elem, add to accum, increment counter, conditional jump to body
	// TODO these can be simplified with the helpers in vector.go
	cond_rhs := env.Block.NewLoad(
		types.I64,
		env.Block.NewGetElementPtr(
			prism.VectorType{Type: typ}.Realise(),
			val.Value,
			I32(0),
			vectorLenOffset))

	ll_body_actual := env.Block.NewLoad(
		types.NewPointer(lltyp),
		env.Block.NewGetElementPtr(
			prism.VectorType{Type: typ}.Realise(),
			val.Value,
			I32(0),
			vectorBodyOffset))

	loopblock := env.CurrentFunction.NewBlock("")
	env.Block.NewBr(loopblock)
	env.Block = loopblock

	// Add to accum
	cur_counter := loopblock.NewLoad(types.I64, counter)

	// Accum <- accum + current element
	loopblock.NewStore(
		env.AgnosticAdd(
			&typ,
			loopblock.NewLoad(lltyp, accum),
			loopblock.NewLoad(
				lltyp,
				loopblock.NewGetElementPtr(
					lltyp,
					ll_body_actual,
					cur_counter))),
		accum)

	// Increment counter
	loopblock.NewStore(
		loopblock.NewAdd(loopblock.NewLoad(types.I64, counter), I64(1)),
		counter)

	cond := loopblock.NewICmp(
		enum.IPredSLT,
		loopblock.NewAdd(cur_counter, I64(1)),
		cond_rhs)

	exitblock := env.CurrentFunction.NewBlock("")
	loopblock.NewCondBr(cond, loopblock, exitblock)
	env.Block = exitblock

	return env.Block.NewLoad(lltyp, accum)
}

func (env *Environment) CompileInlineProduct(val Value) value.Value {
	typ := val.Type.(prism.VectorType).Type
	lltyp := typ.Realise()

	counter := env.Block.NewAlloca(types.I64)
	env.Block.NewStore(I64(0), counter)

	accum := env.Block.NewAlloca(lltyp)
	env.Block.NewStore(env.Number(&typ, 1), accum)

	// Body
	// Get elem, add to accum, increment counter, conditional jump to body

	cond_rhs := env.Block.NewLoad(
		types.I64,
		env.Block.NewGetElementPtr(
			prism.VectorType{Type: typ}.Realise(),
			val.Value,
			I32(0),
			vectorLenOffset))

	ll_body_actual := env.Block.NewLoad(
		types.NewPointer(lltyp),
		env.Block.NewGetElementPtr(
			prism.VectorType{Type: typ}.Realise(),
			val.Value,
			I32(0),
			vectorBodyOffset))

	loopblock := env.CurrentFunction.NewBlock("")
	env.Block.NewBr(loopblock)
	env.Block = loopblock

	// Add to accum
	cur_counter := loopblock.NewLoad(types.I64, counter)

	// Accum <- accum * current element
	loopblock.NewStore(
		env.AgnosticMult(
			&typ,
			loopblock.NewLoad(lltyp, accum),
			loopblock.NewLoad(
				lltyp,
				loopblock.NewGetElementPtr(
					lltyp,
					ll_body_actual,
					cur_counter))),
		accum)

	cond := loopblock.NewICmp(
		enum.IPredSLT,
		loopblock.NewAdd(cur_counter, I64(1)),
		cond_rhs)

	// Increment counter
	loopblock.NewStore(
		loopblock.NewAdd(loopblock.NewLoad(types.I64, counter), I64(1)),
		counter)

	exitblock := env.CurrentFunction.NewBlock("")
	loopblock.NewCondBr(cond, loopblock, exitblock)
	env.Block = exitblock

	return env.Block.NewLoad(lltyp, accum)

}
