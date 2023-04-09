package types

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gridiron "github.com/gridiron-zone/gridiron/types"
)

// DefaultGenesis gets the raw genesis raw message for testing
func DefaultGenesis() *stakingtypes.GenesisState {
	params := stakingtypes.DefaultParams()
	params.BondDenom = gridiron.BaseDenom
	return &stakingtypes.GenesisState{
		Params: params,
	}
}
