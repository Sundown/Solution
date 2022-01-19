package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileInlineMap(fn, vec prism.Expression) value.Value {
	// Return type of function to be mapped
	f_returns := fn.Type()

	// The vector in LLVM
	llvec := env.CompileExpression(&vec)

	head_type := vec.Type().Realise()
	elm_type := vec.Type().(prism.VectorType).Realise()

	to_head_type := prism.VectorType{Type: f_returns}.Realise()
	to_elm_type := f_returns.Realise()

	should_store := true
	if f_returns.Kind() == prism.VoidType.ID {
		should_store = false
	}

	leng := env.ReadVectorLength(Value{llvec, vec.Type()})

	var head *ir.InstAlloca
	var body *ir.InstBitCast

	if should_store {
		cap := env.ReadVectorCapacity(Value{llvec, vec.Type()})
		head = env.Block.NewAlloca(to_head_type)

		// Copy length
		env.Block.NewStore(leng, env.Block.NewGetElementPtr(to_head_type, head, I32(0), vectorLenOffset))

		// Copy capacity
		env.Block.NewStore(cap, env.Block.NewGetElementPtr(to_head_type, head, I32(0), vectorCapOffset))

		// Allocate a body of capacity * element width, and cast to element type
		body = env.Block.NewBitCast(
			env.Block.NewCall(env.GetCalloc(),
				I32(f_returns.Width()),              // Byte size of elements
				env.Block.NewTrunc(cap, types.I32)), // How much memory to alloc
			types.NewPointer(to_elm_type)) // Cast alloc'd memory to typ
	}

	// --- Loop body ---
	vec_body := env.Block.NewLoad(
		types.NewPointer(elm_type),
		env.Block.NewGetElementPtr(head_type, llvec, I32(0), vectorBodyOffset))

	counter := env.Block.NewAlloca(types.I64)
	env.Block.NewStore(I64(0), counter)

	// Body
	// Get elem, add to accum, increment counter, conditional jump to body
	loopblock := env.CurrentFunction.NewBlock("")
	env.Block.NewBr(loopblock)
	env.Block = loopblock
	// Add to accum
	cur_counter := loopblock.NewLoad(types.I64, counter)

	var cur_elm value.Value = loopblock.NewGetElementPtr(elm_type, vec_body, cur_counter)

	if _, ok := vec.Type().(prism.VectorType).Type.(prism.AtomicType); ok {
		cur_elm = loopblock.NewLoad(elm_type, cur_elm)
	}

	var call value.Value

	/*if arg.Morpheme.Tuple[0].Morpheme.Function.Special {
		call = env.GetSpecialCallable(fn.(prism.Function).Ident())(vec.Type().(prism.VectorType), cur_elm)
	} else {*/
	call = loopblock.NewCall(
		env.CompileExpression(&fn),
		cur_elm)

	//}

	if should_store {
		loopblock.NewStore(
			call,
			loopblock.NewGetElementPtr(to_elm_type, body, cur_counter))
	}

	// Increment counter
	incr := loopblock.NewAdd(cur_counter, I64(1))

	loopblock.NewStore(incr, counter)

	cond := loopblock.NewICmp(enum.IPredSLT, incr, leng)

	exitblock := env.CurrentFunction.NewBlock("")
	loopblock.NewCondBr(cond, loopblock, exitblock)
	env.Block = exitblock

	if should_store {
		env.Block.NewStore(body,
			env.Block.NewGetElementPtr(to_head_type, head, I32(0), vectorBodyOffset))
		return head
	} else {
		return nil
	}
}
