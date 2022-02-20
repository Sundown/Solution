package prism

import "github.com/llir/llvm/ir/types"

func (a AtomicType) Realise() types.Type {
	return a.Actual
}

func (v VectorType) Realise() types.Type {
	return types.NewStruct(
		types.I32, types.I32,
		types.NewPointer(v.Type.Realise()))
}

func (s StructType) Realise() types.Type {
	acc := []types.Type{}
	for _, v := range s.FieldTypes {
		acc = append(acc, v.Realise())
	}

	return types.NewStruct(acc...)
}

func (s SumType) Realise() types.Type {
	panic("Impossible")
}

func (s GenericType) Realise() types.Type {
	panic("Impossible")
}
func (f DyadicFunction) LLVMise() string {
	return f.Name.Package + "::" + f.Name.Name + "_" + f.AlphaType.String() + "," + f.OmegaType.String() + "->" + f.Returns.String()
}

func (f MonadicFunction) LLVMise() string {
	return f.Name.Package + "::" + f.Name.Name + "_" + f.OmegaType.String() + "->" + f.Returns.String()
}
