package datacap

import (
	"github.com/CluEleSsUK/go-state-types/abi"
	"github.com/CluEleSsUK/go-state-types/big"
	"github.com/CluEleSsUK/go-state-types/builtin/v9/util/adt"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

type State struct {
	Governor address.Address
	Token    TokenState
}

type TokenState struct {
	Supply       abi.TokenAmount
	Balances     cid.Cid // HAMT abi.ActorID[abi.TokenAmount]
	Allowances   cid.Cid // HAMT abi.ActorID[abi.ActorID[abi.TokenAmount]]
	HamtBitWidth uint64  // uint32 in builtin-actors. uint64 here to satisfy cbor-gen
}

func ConstructState(store adt.Store, governor address.Address, bitWidth uint64) (*State, error) {
	emptyMapCid, err := adt.StoreEmptyMap(store, int(bitWidth))
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty map: %w", err)
	}

	return &State{
		Governor: governor,
		Token: TokenState{
			Supply:       big.Zero(),
			Balances:     emptyMapCid,
			Allowances:   emptyMapCid,
			HamtBitWidth: bitWidth,
		},
	}, nil
}
