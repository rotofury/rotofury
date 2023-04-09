package oracle

import (
	"time"

	gridiron "github.com/gridiron-zone/gridiron/types"
	"github.com/gridiron-zone/gridiron/x/oracle/keeper"
	"github.com/gridiron-zone/gridiron/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	params := k.GetParams(ctx)
	if gridiron.IsPeriodLastBlock(ctx, params.VotePeriod) {
		stakingKeeper := k.StakingKeeper()

		// Build claim map over all validators in active set
		validatorClaimMap := make(map[string]types.Claim)

		maxValidators := stakingKeeper.MaxValidators(ctx)
		iterator := stakingKeeper.ValidatorsPowerStoreIterator(ctx)
		defer iterator.Close()

		powerReduction := stakingKeeper.PowerReduction(ctx)

		i := 0
		for ; iterator.Valid() && i < int(maxValidators); iterator.Next() {
			validator := stakingKeeper.Validator(ctx, iterator.Value())

			// Exclude not bonded validator
			if validator.IsBonded() {
				valAddr := validator.GetOperator()
				validatorClaimMap[valAddr.String()] = types.NewClaim(validator.GetConsensusPower(powerReduction), 0, 0, valAddr)
				i++
			}
		}

		// Denom map
		voteTargets := make(map[string]struct{})
		k.IterateVoteTargets(ctx, func(denom string) bool {
			voteTargets[denom] = struct{}{}
			return false
		})

		// Clear all exchange rates
		k.IterateExchangeRates(ctx, func(denom string, _ sdk.Dec) (stop bool) {
			k.DeleteExchangeRate(ctx, denom)
			return false
		})

		// Organize votes to ballot by denom
		// NOTE: **Filter out inactive or jailed validators**
		// NOTE: **Make abstain votes to have zero vote power**
		voteMap := k.OrganizeBallotByDenom(ctx, validatorClaimMap)
		ctx.Logger().Debug("organized ballot by denom", "voteMap", voteMap)

		if referenceGridb := PickReferenceGridb(ctx, k, voteTargets, voteMap); referenceGridb != "" {
			// make voteMap of Reference Grid to calculate cross exchange rates
			ballotRM := voteMap[referenceGridb]
			voteMapRM := ballotRM.ToMap()

			var exchangeRateRM sdk.Dec

			exchangeRateRM = ballotRM.WeightedMedianWithAssertion()

			// Iterate through ballots and update exchange rates; drop if not enough votes have been achieved.
			for denom, ballot := range voteMap {

				// Convert ballot to cross exchange rates
				if denom != referenceGridb {
					ballot = ballot.ToCrossRateWithSort(voteMapRM)
				}

				// Get weighted median of cross exchange rates
				exchangeRate := Tally(ctx, ballot, params.RewardBand, validatorClaimMap)

				// Transform into the original form {denom}/uUSD
				if denom != referenceGridb {
					if exchangeRate.IsZero() {
						k.Logger(ctx).Error("invalid tallied cross exchange rate", "denom", denom, "exchangeRate", exchangeRate)
						// Do not set exchange rate
						continue
					}
					exchangeRate = exchangeRateRM.Quo(exchangeRate)
				}

				// Set the exchange rate, emit ABCI event
				k.SetExchangeRate(ctx, denom, exchangeRate)
			}
		}

		// ---------------------------
		// Do miss counting & slashing
		voteTargetsLen := len(voteTargets)
		for _, claim := range validatorClaimMap {
			// Skip abstain & valid voters
			if int(claim.WinCount) == voteTargetsLen {
				continue
			}

			// Increase miss counter
			k.SetMissCounter(ctx, claim.Recipient, k.GetMissCounter(ctx, claim.Recipient)+1)
		}

		// Distribute rewards to ballot winners
		k.RewardBallotWinners(
			ctx,
			(int64)(params.VotePeriod),
			(int64)(params.RewardDistributionWindow),
			voteTargets,
			validatorClaimMap,
		)

		// Clear the ballot
		k.ClearBallots(ctx, params.VotePeriod)
	}

	// Do slash who did miss voting over threshold and
	// reset miss counters of all validators
	// at the last block of slash window
	if gridiron.IsPeriodLastBlock(ctx, params.SlashWindow) {
		k.SlashAndResetMissCounters(ctx)
	}

	return
}
