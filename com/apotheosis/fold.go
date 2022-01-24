package apotheosis

// Need tuple element access before this is useful, ignore for now
/* func (env *Environment) CompileInlineFoldl(app *prism.DApplication) value.Value {
	fn := app.Alpha.Morpheme.Tuple[0]
	llfn := env.CompileExpression(fn)
	typ := app.ArgumentAlpha.TypeOf.Vector

	lltyp := typ.AsLLType()

	vec := app.ArgumentAlpha

	llvec := env.CompileExpression(vec)

	counter := env.Block.NewAlloca(types.I32)
	env.Block.NewStore(I32(0), counter)

	accum := env.Block.NewAlloca(lltyp)
	env.Block.NewStore(env.Number(typ, 1), accum)

	// Body
	// Get elem, add to accum, increment counter, conditional jump to body

	cond_rhs := env.Block.NewLoad(
		types.I32,
		env.Block.NewGetElementPtr(
			typ.AsVector().AsLLType(),
			llvec,
			I32(0),
			vectorLenOffset))

	// This is a bit messy
	loopblock := env.CurrentFunction.NewBlock("")
	env.Block.NewBr(loopblock)
	env.Block = loopblock
	// ---

	// Add to accum
	cur_counter := loopblock.NewLoad(types.I32, counter)

	// Accum <- accum * current element
	ll_tuple := env.Block.NewAlloca(fn.TypeOf.AsLLType())

	left := accum
	right := loopblock.NewGetElementPtr(
		lltyp,
		env.Block.NewLoad(
			types.NewPointer(lltyp),
			env.Block.NewGetElementPtr(
				typ.AsVector().AsLLType(),
				llvec,
				I32(0),
				vectorBodyOffset)),
		cur_counter)

	env.Block.NewStore(left, env.GEP(ll_tuple, I32(0), I32(0)))
	env.Block.NewStore(right, env.GEP(ll_tuple, I32(0), I32(1)))

	loopblock.NewStore(env.Block.NewCall(llfn, ll_tuple), accum)

	cond := loopblock.NewICmp(
		enum.IPredSLT,
		loopblock.NewAdd(cur_counter, I32(1)),
		cond_rhs)

	// Increment counter
	loopblock.NewStore(
		loopblock.NewAdd(loopblock.NewLoad(types.I32, counter), I32(1)),
		counter)

	exitblock := env.CurrentFunction.NewBlock("")
	loopblock.NewCondBr(cond, loopblock, exitblock)
	env.Block = exitblock

	return env.Block.NewLoad(lltyp, accum)
}
*/
