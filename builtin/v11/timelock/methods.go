package timelock

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)},
	2: {"Decrypt", *new(func(params *DecryptParams) *DecryptParams)},
	builtin.MustGenerateFRCMethodNum("InvokeEVM"): {"InvokeContract", *new(func(bytes *abi.CborBytes) *abi.CborBytes)},
}
