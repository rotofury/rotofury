package vesting

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gridiron-zone/gridiron/x/vesting/keeper"
	"github.com/gridiron-zone/gridiron/x/vesting/types"
)

// InitGenesis initializes the vesting module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	allocAddresses := k.GetAllocationAddresses(ctx)
	if len(genState.AllocationAddresses.StrategicReserveCustodianAddr) != 0 {
		newSrca, err := sdk.AccAddressFromBech32(genState.AllocationAddresses.StrategicReserveCustodianAddr)
		if err != nil {
			panic(err)
		}
		allocAddresses.StrategicReserveCustodianAddr = newSrca.String()
	}
	if len(genState.AllocationAddresses.TeamVestingAddr) != 0 {
		newTva, err := sdk.AccAddressFromBech32(genState.AllocationAddresses.TeamVestingAddr)
		if err != nil {
			panic(err)
		}
		allocAddresses.TeamVestingAddr = newTva.String()
	}
	k.SetAllocationAddresses(ctx, allocAddresses)

	if ctx.BlockHeight() <= 1 {
		k.AllocateAtGenesis(ctx, genState)
	}
}

// ExportGenesis returns the vesting module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Params = k.GetParams(ctx)
	genesis.AllocationAddresses = k.GetAllocationAddresses(ctx)

	return genesis
}
