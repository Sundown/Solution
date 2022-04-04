package apotheosis

/*
func (env *Environment) compileTuple(tuple *prism.Morpheme) value.Value {
	ll_tuple := env.Block.NewAlloca(tuple.TypeOf.AsLLType())

	for index, expr := range tuple.Tuple {
		val := env.compileExpression(expr)

		if expr.TypeOf.Atomic == nil {
			val = env.Block.NewLoad(expr.TypeOf.AsLLType(), val)
		}

		env.Block.NewStore(val, env.gep(ll_tuple, i32(0), i32(int64(index))))
	}

	return ll_tuple
}

func (env *Environment) TupleGet(typ *prism.Type, real value.Value, index int) value.Value {
	if len(typ.Tuple) < index {
		prism.Panic(prism.CT_OOB, index, typ.String(), len(typ.Tuple))
	}

	if typ.Tuple == nil {
		prism.Panic(prism.CT_Unexpected, prism.Yellow("tuple"), prism.Yellow(typ.String()))
	}

	ptr := env.Block.NewGetElementPtr(
		typ.AsLLType(), real,
		i32(0), i32(int64(index)))

	if typ.Tuple[index].Atomic == nil {
		return ptr
	} else {
		return env.Block.NewLoad(typ.Tuple[index].AsLLType(), ptr)
	}
}
*/
