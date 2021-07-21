package compiler

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func I64(v int64) constant.Constant {
	return constant.NewInt(types.I64, v)

}

func I32(v int64) constant.Constant {
	return constant.NewInt(types.I32, int64(int32(v)))

}
