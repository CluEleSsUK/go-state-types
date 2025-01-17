package verifreg

import (
	"github.com/CluEleSsUK/go-state-types/builtin"
)

const EndOfLifeClaimDropPeriod = 30 * builtin.EpochsInDay

const MaximumVerifiedAllocationExpiration = 60 * builtin.EpochsInDay

const MinimumVerifiedAllocationTerm = 180 * builtin.EpochsInDay

const MaximumVerifiedAllocationTerm = 5 * builtin.EpochsInYear

const NoAllocationID = AllocationId(0)
