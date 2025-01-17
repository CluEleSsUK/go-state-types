package v8

import (
	"bytes"

	"github.com/CluEleSsUK/go-state-types/manifest"
	"github.com/ipfs/go-cid"

	"github.com/CluEleSsUK/go-state-types/builtin/v8/account"
	"github.com/CluEleSsUK/go-state-types/builtin/v8/cron"
	init_ "github.com/CluEleSsUK/go-state-types/builtin/v8/init"
	"github.com/CluEleSsUK/go-state-types/builtin/v8/market"
	"github.com/CluEleSsUK/go-state-types/builtin/v8/miner"
	"github.com/CluEleSsUK/go-state-types/builtin/v8/multisig"
	"github.com/CluEleSsUK/go-state-types/builtin/v8/paych"
	"github.com/CluEleSsUK/go-state-types/builtin/v8/power"
	"github.com/CluEleSsUK/go-state-types/builtin/v8/reward"
	"github.com/CluEleSsUK/go-state-types/builtin/v8/verifreg"

	"github.com/CluEleSsUK/go-state-types/builtin"

	"github.com/CluEleSsUK/go-state-types/abi"
	"github.com/CluEleSsUK/go-state-types/big"
	"github.com/filecoin-project/go-address"
	"golang.org/x/xerrors"
)

// Within this code, Go errors are not expected, but are often converted to messages so that execution
// can continue to find more errors rather than fail with no insight.
// Only errors thar are particularly troublesome to recover from should propagate as Go errors.
func CheckStateInvariants(tree *builtin.ActorTree, priorEpoch abi.ChainEpoch, actorCodes map[string]cid.Cid) (*builtin.MessageAccumulator, error) {
	acc := &builtin.MessageAccumulator{}
	totalFIl := big.Zero()
	var initSummary *init_.StateSummary
	var cronSummary *cron.StateSummary
	var verifregSummary *verifreg.StateSummary
	var marketSummary *market.StateSummary
	var rewardSummary *reward.StateSummary
	var accountSummaries []*account.StateSummary
	var powerSummary *power.StateSummary
	var paychSummaries []*paych.StateSummary
	var multisigSummaries []*multisig.StateSummary
	minerSummaries := make(map[address.Address]*miner.StateSummary)

	if err := tree.ForEachV4(func(key address.Address, actor *builtin.ActorV4) error {
		acc := acc.WithPrefix("%v ", key) // Intentional shadow
		if key.Protocol() != address.ID {
			acc.Addf("unexpected address protocol in state tree root: %v", key)
		}
		totalFIl = big.Add(totalFIl, actor.Balance)

		switch actor.Code {
		case actorCodes[manifest.SystemKey]:

		case actorCodes[manifest.InitKey]:
			var st init_.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := init_.CheckStateInvariants(&st, tree.Store)
			acc.WithPrefix("init: ").AddAll(msgs)
			initSummary = summary
		case actorCodes[manifest.CronKey]:
			var st cron.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := cron.CheckStateInvariants(&st, tree.Store)
			acc.WithPrefix("cron: ").AddAll(msgs)
			cronSummary = summary
		case actorCodes[manifest.AccountKey]:
			var st account.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := account.CheckStateInvariants(&st, key)
			acc.WithPrefix("account: ").AddAll(msgs)
			accountSummaries = append(accountSummaries, summary)
		case actorCodes[manifest.PowerKey]:
			var st power.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := power.CheckStateInvariants(&st, tree.Store)
			acc.WithPrefix("power: ").AddAll(msgs)
			powerSummary = summary
		case actorCodes[manifest.MinerKey]:
			var st miner.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := miner.CheckStateInvariants(&st, tree.Store, actor.Balance)
			acc.WithPrefix("miner: ").AddAll(msgs)
			minerSummaries[key] = summary
		case actorCodes[manifest.MarketKey]:
			var st market.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := market.CheckStateInvariants(&st, tree.Store, actor.Balance, priorEpoch)
			acc.WithPrefix("market: ").AddAll(msgs)
			marketSummary = summary
		case actorCodes[manifest.PaychKey]:
			var st paych.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := paych.CheckStateInvariants(&st, tree.Store, actor.Balance)
			acc.WithPrefix("paych: ").AddAll(msgs)
			paychSummaries = append(paychSummaries, summary)
		case actorCodes[manifest.MultisigKey]:
			var st multisig.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := multisig.CheckStateInvariants(&st, tree.Store)
			acc.WithPrefix("multisig: ").AddAll(msgs)
			multisigSummaries = append(multisigSummaries, summary)
		case actorCodes[manifest.RewardKey]:
			var st reward.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := reward.CheckStateInvariants(&st, tree.Store, priorEpoch, actor.Balance)
			acc.WithPrefix("reward: ").AddAll(msgs)
			rewardSummary = summary
		case actorCodes[manifest.VerifregKey]:
			var st verifreg.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			summary, msgs := verifreg.CheckStateInvariants(&st, tree.Store)
			acc.WithPrefix("verifreg: ").AddAll(msgs)
			verifregSummary = summary
		default:
			return xerrors.Errorf("unexpected actor code CID %v for address %v", actor.Code, key)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	//
	// Perform cross-actor checks from state summaries here.
	//

	CheckMinersAgainstPower(acc, minerSummaries, powerSummary)
	CheckDealStatesAgainstSectors(acc, minerSummaries, marketSummary)

	_ = initSummary
	_ = verifregSummary
	_ = cronSummary
	_ = marketSummary
	_ = rewardSummary

	if !totalFIl.Equals(builtin.TotalFilecoin) {
		acc.Addf("total token balance is %v, expected %v", totalFIl, builtin.TotalFilecoin)
	}

	return acc, nil
}

func CheckMinersAgainstPower(acc *builtin.MessageAccumulator, minerSummaries map[address.Address]*miner.StateSummary, powerSummary *power.StateSummary) {
	for addr, minerSummary := range minerSummaries { // nolint:nomaprange
		// check claim
		claim, ok := powerSummary.Claims[addr]
		acc.Require(ok, "miner %v has no power claim", addr)
		if ok {
			claimPower := miner.NewPowerPair(claim.RawBytePower, claim.QualityAdjPower)
			acc.Require(minerSummary.ActivePower.Equals(claimPower),
				"miner %v computed active power %v does not match claim %v", addr, minerSummary.ActivePower, claimPower)
			acc.Require(minerSummary.WindowPoStProofType == claim.WindowPoStProofType,
				"miner seal proof type %d does not match claim proof type %d", minerSummary.WindowPoStProofType, claim.WindowPoStProofType)
		}

		// check crons
		crons, ok := powerSummary.Crons[addr]
		if !ok { // with deferred and discontinued crons it is normal for a miner actor to have no cron events
			continue
		}

		var payload miner.CronEventPayload
		var provingPeriodCron *power.MinerCronEvent
		for _, event := range crons {
			err := payload.UnmarshalCBOR(bytes.NewReader(event.Payload))
			acc.Require(err == nil, "miner %v registered cron at epoch %d with wrong or corrupt payload",
				addr, event.Epoch)
			acc.Require(payload.EventType == miner.CronEventProcessEarlyTerminations || payload.EventType == miner.CronEventProvingDeadline,
				"miner %v has unexpected cron event type %v", addr, payload.EventType)

			if payload.EventType == miner.CronEventProvingDeadline {
				if provingPeriodCron != nil {
					acc.Require(false, "miner %v has duplicate proving period crons at epoch %d and %d",
						addr, provingPeriodCron.Epoch, event.Epoch)
				}
				provingPeriodCron = &event
			}
		}
		hasProvingPeriodCron := provingPeriodCron != nil
		acc.Require(hasProvingPeriodCron == minerSummary.DeadlineCronActive, "miner %v has invalid DeadlineCronActive (%t) for hasProvingPeriodCron status (%t)",
			addr, minerSummary.DeadlineCronActive, hasProvingPeriodCron)

		acc.Require(provingPeriodCron != nil, "miner %v has no proving period cron", addr)
	}
}

func CheckDealStatesAgainstSectors(acc *builtin.MessageAccumulator, minerSummaries map[address.Address]*miner.StateSummary, marketSummary *market.StateSummary) {
	// Check that all active deals are included within a non-terminated sector.
	// We cannot check that all deals referenced within a sector are in the market, because deals
	// can be terminated independently of the sector in which they are included.
	for dealID, deal := range marketSummary.Deals { // nolint:nomaprange
		if deal.SectorStartEpoch == abi.ChainEpoch(-1) {
			// deal hasn't been activated yet, make no assertions about sector state
			continue
		}

		minerSummary, found := minerSummaries[deal.Provider]
		if !found {
			acc.Addf("provider %v for deal %d not found among miners", deal.Provider, dealID)
			continue
		}

		sectorDeal, found := minerSummary.Deals[dealID]
		if !found {
			acc.Require(deal.SlashEpoch >= 0, "un-slashed deal %d not referenced in active sectors of miner %v", dealID, deal.Provider)
			continue
		}

		acc.Require(deal.SectorStartEpoch == sectorDeal.SectorStart,
			"deal state start %d does not match sector start %d for miner %v",
			deal.SectorStartEpoch, sectorDeal.SectorStart, deal.Provider)

		acc.Require(deal.SectorStartEpoch <= sectorDeal.SectorExpiration,
			"deal state start %d activated after sector expiration %d for miner %v",
			deal.SectorStartEpoch, sectorDeal.SectorExpiration, deal.Provider)

		acc.Require(deal.LastUpdatedEpoch <= sectorDeal.SectorExpiration,
			"deal state update at %d after sector expiration %d for miner %v",
			deal.LastUpdatedEpoch, sectorDeal.SectorExpiration, deal.Provider)

		acc.Require(deal.SlashEpoch <= sectorDeal.SectorExpiration,
			"deal state slashed at %d after sector expiration %d for miner %v",
			deal.SlashEpoch, sectorDeal.SectorExpiration, deal.Provider)
	}
}
