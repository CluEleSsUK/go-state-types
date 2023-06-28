package power

import (
	"github.com/CluEleSsUK/go-state-types/abi"
	"github.com/CluEleSsUK/go-state-types/builtin/v8/util/smoothing"
	addr "github.com/filecoin-project/go-address"
)

// Storage miner actor constructor params are defined here so the power actor can send them to the init actor
// to instantiate miners.
// Changed since v2:
// - Seal proof type replaced with PoSt proof type
type MinerConstructorParams struct {
	OwnerAddr           addr.Address
	WorkerAddr          addr.Address
	ControlAddrs        []addr.Address
	WindowPoStProofType abi.RegisteredPoStProof
	PeerId              abi.PeerID
	Multiaddrs          []abi.Multiaddrs
}

type CreateMinerParams struct {
	Owner               addr.Address
	Worker              addr.Address
	WindowPoStProofType abi.RegisteredPoStProof
	Peer                abi.PeerID
	Multiaddrs          []abi.Multiaddrs
}

type CreateMinerReturn struct {
	IDAddress     addr.Address // The canonical ID-based address for the actor.
	RobustAddress addr.Address // A more expensive but re-org-safe address for the newly created actor.
}

type UpdateClaimedPowerParams struct {
	RawByteDelta         abi.StoragePower
	QualityAdjustedDelta abi.StoragePower
}

type EnrollCronEventParams struct {
	EventEpoch abi.ChainEpoch
	Payload    []byte
}

type CurrentTotalPowerReturn struct {
	RawBytePower            abi.StoragePower
	QualityAdjPower         abi.StoragePower
	PledgeCollateral        abi.TokenAmount
	QualityAdjPowerSmoothed smoothing.FilterEstimate
}
