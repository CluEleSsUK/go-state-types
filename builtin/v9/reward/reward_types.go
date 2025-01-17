package reward

import (
	"github.com/CluEleSsUK/go-state-types/abi"
	"github.com/CluEleSsUK/go-state-types/builtin/v9/util/smoothing"
	"github.com/filecoin-project/go-address"
)

type AwardBlockRewardParams struct {
	Miner     address.Address
	Penalty   abi.TokenAmount // penalty for including bad messages in a block, >= 0
	GasReward abi.TokenAmount // gas reward from all gas fees in a block, >= 0
	WinCount  int64           // number of reward units won, > 0
}

type ThisEpochRewardReturn struct {
	ThisEpochRewardSmoothed smoothing.FilterEstimate
	ThisEpochBaselinePower  abi.StoragePower
}
