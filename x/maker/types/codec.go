package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgMintBySwap{}, "gridiron/MsgMintBySwap", nil)
	cdc.RegisterConcrete(&MsgBurnBySwap{}, "gridiron/MsgBurnBySwap", nil)
	cdc.RegisterConcrete(&MsgBuyBacking{}, "gridiron/MsgBuyBacking", nil)
	cdc.RegisterConcrete(&MsgSellBacking{}, "gridiron/MsgSellBacking", nil)
	cdc.RegisterConcrete(&MsgMintByCollateral{}, "gridiron/MsgMintByCollateral", nil)
	cdc.RegisterConcrete(&MsgBurnByCollateral{}, "gridiron/MsgBurnByCollateral", nil)
	cdc.RegisterConcrete(&MsgDepositCollateral{}, "gridiron/MsgDepositCollateral", nil)
	cdc.RegisterConcrete(&MsgRedeemCollateral{}, "gridiron/MsgRedeemCollateral", nil)
	cdc.RegisterConcrete(&MsgLiquidateCollateral{}, "gridiron/MsgLiquidateCollateral", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&RegisterBackingProposal{},
		&RegisterCollateralProposal{},
		&SetBackingRiskParamsProposal{},
		&SetCollateralRiskParamsProposal{},
		&BatchSetBackingRiskParamsProposal{},
		&BatchSetCollateralRiskParamsProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterCodec(Amino)
}
