package types

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gridiron "github.com/gridiron-zone/gridiron/types"
)

const (
	EmptyVeID = 0
	FirstVeID = 1
	MaxVeID   = math.MaxUint64 - 1

	// 4 years, i.e., 209 weeks
	MaxLockTimeWeeks = gridiron.DaysPer4Years/7 + 1
	MaxLockTime      = MaxLockTimeWeeks * gridiron.SecondsPerWeek

	MaxUnixTime = math.MaxInt64

	// Regulated period for ve locking time
	RegulatedPeriod = gridiron.SecondsPerWeek

	EmptyEpoch = 0
	FirstEpoch = 1
)

var (
	// Emission amount are halved every 4 years (almost 209 weeks).
	// For geometric sequence of every 4 years,
	// a * (r ^ n) = a * 0.5 where n = 209
	// so that <emission ratio per week> ^ 209 = 0.5
	EmissionRatio, _ = sdk.NewDecFromStr("0.9966889998035777")

	// Minimum circulating rate allowed for calculating emission
	MinEmissionCirculating = sdk.NewDecWithPrec(1, 10)
)
