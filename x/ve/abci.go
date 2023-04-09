package ve

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gridiron-zone/gridiron/x/ve/keeper"
	"github.com/gridiron-zone/gridiron/x/ve/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	k.RegulateCheckpoint(ctx)
}
