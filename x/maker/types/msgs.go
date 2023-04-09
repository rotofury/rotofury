package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gridiron "github.com/gridiron-zone/gridiron/types"
)

const (
	TypeMsgMintBySwap          = "mint_by_swap"
	TypeMsgBurnBySwap          = "burn_by_swap"
	TypeMsgMintByCollateral    = "mint_by_collateral"
	TypeMsgBurnByCollateral    = "burn_by_collateral"
	TypeMsgDepositCollateral   = "deposit_collateral"
	TypeMsgRedeemCollateral    = "redeem_collateral"
	TypeMsgBuyBacking          = "buy_backing"
	TypeMsgSellBacking         = "sell_backing"
	TypeMsgLiquidateCollateral = "liquidate_collateral"
)

var (
	_ sdk.Msg = &MsgMintBySwap{}
	_ sdk.Msg = &MsgBurnBySwap{}
	_ sdk.Msg = &MsgMintByCollateral{}
	_ sdk.Msg = &MsgBurnByCollateral{}
	_ sdk.Msg = &MsgDepositCollateral{}
	_ sdk.Msg = &MsgRedeemCollateral{}
	_ sdk.Msg = &MsgBuyBacking{}
	_ sdk.Msg = &MsgSellBacking{}
	_ sdk.Msg = &MsgLiquidateCollateral{}
)

// Route implements sdk.Msg
func (m *MsgMintBySwap) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgMintBySwap) Type() string { return TypeMsgMintBySwap }

// GetSignBytes implements sdk.Msg
func (m *MsgMintBySwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgMintBySwap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.To) > 0 {
		_, err = sdk.AccAddressFromBech32(m.To)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if m.MintOutMin.Denom != gridiron.MicroUSMDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.MintOutMin.Denom)
	}
	if m.IronInMax.Denom != gridiron.AttoIronDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.IronInMax.Denom)
	}
	if !m.MintOutMin.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "mint_out_min must be positive")
	}
	if m.BackingInMax.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "backing_in_max must be positive or zero")
	}
	if m.IronInMax.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "iron_in_max must be positive or zero")
	}
	if m.BackingInMax.Amount.IsZero() && m.IronInMax.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "backing_in_max and iron_in_max must not be both zero")
	}
	if m.FullBacking && m.IronInMax.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "iron_in_max must be zero when full_backing is true")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgMintBySwap) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgBurnBySwap) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgBurnBySwap) Type() string { return TypeMsgBurnBySwap }

// GetSignBytes implements sdk.Msg
func (m *MsgBurnBySwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgBurnBySwap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.To) > 0 {
		_, err = sdk.AccAddressFromBech32(m.To)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if m.BurnIn.Denom != gridiron.MicroUSMDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.BurnIn.Denom)
	}
	if m.IronOutMin.Denom != gridiron.AttoIronDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.IronOutMin.Denom)
	}
	if !m.BurnIn.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.BurnIn.String())
	}
	if m.BackingOutMin.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.BackingOutMin.String())
	}
	if m.IronOutMin.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.IronOutMin.String())
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgBurnBySwap) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgBuyBacking) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgBuyBacking) Type() string { return TypeMsgBuyBacking }

// GetSignBytes implements sdk.Msg
func (m *MsgBuyBacking) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgBuyBacking) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.To) > 0 {
		_, err = sdk.AccAddressFromBech32(m.To)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if m.IronIn.Denom != gridiron.AttoIronDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.IronIn.Denom)
	}
	if !m.IronIn.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.IronIn.String())
	}
	if m.BackingOutMin.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.BackingOutMin.String())
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgBuyBacking) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgSellBacking) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgSellBacking) Type() string { return TypeMsgSellBacking }

// GetSignBytes implements sdk.Msg
func (m *MsgSellBacking) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgSellBacking) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.To) > 0 {
		_, err = sdk.AccAddressFromBech32(m.To)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if m.IronOutMin.Denom != gridiron.AttoIronDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.IronOutMin.Denom)
	}
	if !m.BackingIn.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.BackingIn.String())
	}
	if m.IronOutMin.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.IronOutMin.String())
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgSellBacking) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgMintByCollateral) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgMintByCollateral) Type() string { return TypeMsgMintByCollateral }

// GetSignBytes implements sdk.Msg
func (m *MsgMintByCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgMintByCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.To) > 0 {
		_, err = sdk.AccAddressFromBech32(m.To)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if m.MintOut.Denom != gridiron.MicroUSMDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.MintOut.Denom)
	}
	if !m.MintOut.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.MintOut.String())
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgMintByCollateral) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgBurnByCollateral) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgBurnByCollateral) Type() string { return TypeMsgBurnByCollateral }

// GetSignBytes implements sdk.Msg
func (m *MsgBurnByCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgBurnByCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if m.RepayInMax.Denom != gridiron.MicroUSMDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.RepayInMax.Denom)
	}
	if !m.RepayInMax.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.RepayInMax.String())
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgBurnByCollateral) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgDepositCollateral) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgDepositCollateral) Type() string { return TypeMsgDepositCollateral }

// GetSignBytes implements sdk.Msg
func (m *MsgDepositCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgDepositCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.To) > 0 {
		_, err = sdk.AccAddressFromBech32(m.To)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if m.CollateralIn.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.CollateralIn.String())
	}
	if m.IronIn.Denom != gridiron.AttoIronDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.IronIn.Denom)
	}
	if m.IronIn.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.IronIn.String())
	}
	if m.CollateralIn.Amount.IsZero() && m.IronIn.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, sdk.NewCoins(m.CollateralIn, m.IronIn).String())
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgDepositCollateral) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgRedeemCollateral) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgRedeemCollateral) Type() string { return TypeMsgRedeemCollateral }

// GetSignBytes implements sdk.Msg
func (m *MsgRedeemCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgRedeemCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.To) > 0 {
		_, err = sdk.AccAddressFromBech32(m.To)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if m.CollateralOut.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.CollateralOut.String())
	}
	if m.IronOut.Denom != gridiron.AttoIronDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.IronOut.Denom)
	}
	if m.IronOut.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.IronOut.String())
	}
	if m.CollateralOut.Amount.IsZero() && m.IronOut.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, sdk.NewCoins(m.CollateralOut, m.IronOut).String())
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgRedeemCollateral) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgLiquidateCollateral) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgLiquidateCollateral) Type() string { return TypeMsgLiquidateCollateral }

// GetSignBytes implements sdk.Msg
func (m *MsgLiquidateCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgLiquidateCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.To) > 0 {
		_, err = sdk.AccAddressFromBech32(m.To)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	_, err = sdk.AccAddressFromBech32(m.Debtor)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid debtor address (%s)", err)
	}
	if !m.Collateral.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Collateral.String())
	}
	if m.RepayInMax.Denom != gridiron.MicroUSMDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin: %s", m.RepayInMax.Denom)
	}
	if !m.RepayInMax.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.RepayInMax.String())
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgLiquidateCollateral) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
