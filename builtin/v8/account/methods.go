package account

import (
	"github.com/filecoin-project/go-address"

	"github.com/CluEleSsUK/go-state-types/abi"
	"github.com/CluEleSsUK/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*address.Address) *abi.EmptyValue)},   // Constructor
	2: {"PubkeyAddress", *new(func(*abi.EmptyValue) *address.Address)}, // PubkeyAddress
}
