package verifreg

import (
	"github.com/CluEleSsUK/go-state-types/abi"
	addr "github.com/filecoin-project/go-address"
)

type AddVerifierParams struct {
	Address   addr.Address
	Allowance DataCap
}

type AddVerifiedClientParams struct {
	Address   addr.Address
	Allowance DataCap
}

type UseBytesParams struct {
	Address  addr.Address     // Address of verified client.
	DealSize abi.StoragePower // Number of bytes to use.
}

type RestoreBytesParams struct {
	Address  addr.Address
	DealSize abi.StoragePower
}

type RemoveDataCapParams struct {
	VerifiedClientToRemove addr.Address
	DataCapAmountToRemove  DataCap
	VerifierRequest1       RemoveDataCapRequest
	VerifierRequest2       RemoveDataCapRequest
}

type RemoveDataCapReturn struct {
	VerifiedClient addr.Address
	DataCapRemoved DataCap
}
