package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gridiron "github.com/gridiron-zone/gridiron/types"
	"github.com/gridiron-zone/gridiron/x/maker/types"
)

func (k Keeper) calculateMintBySwapIn(
	ctx sdk.Context,
	mintOut sdk.Coin,
	backingDenom string,
	fullBacking bool,
) (
	backingIn sdk.Coin,
	ironIn sdk.Coin,
	mintFee sdk.Coin,
	err error,
) {
	backingIn = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	ironIn = sdk.NewCoin(gridiron.AttoIronDenom, sdk.ZeroInt())
	mintFee = sdk.NewCoin(gridiron.MicroUSMDenom, sdk.ZeroInt())

	err = k.checkMintPriceLowerBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	ironPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.AttoIronDenom)
	if err != nil {
		return
	}

	mintFee = computeFee(mintOut, backingParams.MintFee)
	mintTotal := mintOut.Add(mintFee)
	mintTotalInUSD := mintTotal.Amount.ToDec().Mul(gridiron.MicroUSMTarget)

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	poolBacking.GridMinted = poolBacking.GridMinted.Add(mintTotal)
	if backingParams.MaxGridMint != nil && poolBacking.GridMinted.Amount.GT(*backingParams.MaxGridMint) {
		err = sdkerrors.Wrapf(types.ErrGridCeiling, "grid over ceiling")
		return
	}

	backingRatio := k.GetBackingRatio(ctx)
	if backingRatio.GTE(sdk.OneDec()) || fullBacking {
		// full/over backing, or user selects full backing
		backingIn.Amount = mintTotalInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if backingRatio.IsZero() {
		// full algorithmic
		ironIn.Amount = mintTotalInUSD.QuoRoundUp(ironPrice).RoundInt()
	} else {
		// fractional
		backingIn.Amount = mintTotalInUSD.Mul(backingRatio).QuoRoundUp(backingPrice).RoundInt()
		ironIn.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoRoundUp(ironPrice).RoundInt()
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrapf(types.ErrBackingCeiling, "backing over ceiling")
		return
	}

	return
}

func (k Keeper) calculateMintBySwapOut(
	ctx sdk.Context,
	backingInMax sdk.Coin,
	ironInMax sdk.Coin,
	fullBacking bool,
) (
	backingIn sdk.Coin,
	ironIn sdk.Coin,
	mintOut sdk.Coin,
	mintFee sdk.Coin,
	err error,
) {
	backingDenom := backingInMax.Denom

	err = k.checkMintPriceLowerBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in uusd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	ironPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.AttoIronDenom)
	if err != nil {
		return
	}

	backingRatio := k.GetBackingRatio(ctx)

	backingInMaxInUSD := backingPrice.MulInt(backingInMax.Amount)
	ironInMaxInUSD := ironPrice.MulInt(ironInMax.Amount)

	mintTotalInUSD := sdk.ZeroDec()
	backingIn = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	ironIn = sdk.NewCoin(gridiron.AttoIronDenom, sdk.ZeroInt())

	if backingRatio.GTE(sdk.OneDec()) || fullBacking {
		// full/over backing, or user selects full backing
		mintTotalInUSD = backingInMaxInUSD
		backingIn.Amount = backingInMax.Amount
	} else if backingRatio.IsZero() {
		// full algorithmic
		mintTotalInUSD = ironInMaxInUSD
		ironIn.Amount = ironInMax.Amount
	} else {
		// fractional
		max1 := backingInMaxInUSD.Quo(backingRatio)
		max2 := ironInMaxInUSD.Quo(sdk.OneDec().Sub(backingRatio))
		if backingInMax.IsPositive() && (ironInMax.IsZero() || max1.LTE(max2)) {
			mintTotalInUSD = max1
			backingIn.Amount = backingInMax.Amount
			ironIn.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoRoundUp(ironPrice).RoundInt()
			if ironInMax.IsPositive() && ironInMax.IsLT(ironIn) {
				ironIn.Amount = ironInMax.Amount
			}
		} else {
			mintTotalInUSD = max2
			ironIn.Amount = ironInMax.Amount
			backingIn.Amount = mintTotalInUSD.Mul(backingRatio).QuoRoundUp(backingPrice).RoundInt()
			if backingInMax.IsPositive() && backingInMax.IsLT(backingIn) {
				backingIn.Amount = backingInMax.Amount
			}
		}
	}

	mintTotal := sdk.NewCoin(gridiron.MicroUSMDenom, mintTotalInUSD.Quo(gridiron.MicroUSMTarget).TruncateInt())

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	poolBacking.GridMinted = poolBacking.GridMinted.AddAmount(mintTotal.Amount)
	if backingParams.MaxGridMint != nil && poolBacking.GridMinted.Amount.GT(*backingParams.MaxGridMint) {
		err = sdkerrors.Wrap(types.ErrGridCeiling, "")
		return
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "")
		return
	}

	mintFee = computeFee(mintTotal, backingParams.MintFee)
	mintOut = mintTotal.Sub(mintFee)
	return
}

func (k Keeper) calculateBurnBySwapIn(
	ctx sdk.Context,
	backingOutMax sdk.Coin,
	ironOutMax sdk.Coin,
) (
	burnIn sdk.Coin,
	backingOut sdk.Coin,
	ironOut sdk.Coin,
	burnFee sdk.Coin,
	err error,
) {
	backingDenom := backingOutMax.Denom

	burnIn = sdk.NewCoin(gridiron.MicroUSMDenom, sdk.ZeroInt())
	backingOut = sdk.NewCoin(backingOutMax.Denom, sdk.ZeroInt())
	ironOut = sdk.NewCoin(gridiron.AttoIronDenom, sdk.ZeroInt())
	burnFee = sdk.NewCoin(gridiron.MicroUSMDenom, sdk.ZeroInt())

	err = k.checkBurnPriceUpperBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	ironPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.AttoIronDenom)
	if err != nil {
		return
	}

	backingOutMaxInUSD := backingPrice.MulInt(backingOutMax.Amount)
	ironOutMaxInUSD := ironPrice.MulInt(ironOutMax.Amount)

	burnActualInUSD := sdk.ZeroDec()
	backingRatio := k.GetBackingRatio(ctx)
	if backingRatio.GTE(sdk.OneDec()) {
		// full/over backing
		burnActualInUSD = backingOutMaxInUSD
		backingOut.Amount = backingOutMax.Amount
	} else if backingRatio.IsZero() {
		// full algorithmic
		burnActualInUSD = ironOutMaxInUSD
		ironOut.Amount = ironOutMax.Amount
	} else {
		// fractional
		burnActualWithBackingInUSD := backingOutMaxInUSD.Quo(backingRatio)
		burnActualWithIronInUSD := ironOutMaxInUSD.Quo(sdk.OneDec().Sub(backingRatio))
		if ironOutMax.IsZero() || (backingOutMax.IsPositive() && burnActualWithBackingInUSD.LT(burnActualWithIronInUSD)) {
			burnActualInUSD = burnActualWithBackingInUSD
			backingOut.Amount = backingOutMax.Amount
			ironOut.Amount = burnActualInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoRoundUp(ironPrice).RoundInt()
		} else {
			burnActualInUSD = burnActualWithIronInUSD
			ironOut.Amount = ironOutMax.Amount
			backingOut.Amount = burnActualInUSD.Mul(backingRatio).QuoRoundUp(backingPrice).RoundInt()
		}
	}

	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)
	if moduleOwnedBacking.IsLT(backingOut) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) < balance(%s)", backingOut, moduleOwnedBacking)
		return
	}

	burnFeeRate := sdk.ZeroDec()
	if backingParams.BurnFee != nil {
		burnFeeRate = *backingParams.BurnFee
	}

	burnInValue := burnActualInUSD.Quo(gridiron.MicroUSMTarget).Quo(sdk.OneDec().Sub(burnFeeRate))
	burnFeeValue := burnInValue.Mul(burnFeeRate)
	burnIn = sdk.NewCoin(gridiron.MicroUSMDenom, burnInValue.RoundInt())
	burnFee = sdk.NewCoin(gridiron.MicroUSMDenom, burnFeeValue.RoundInt())
	return
}

func (k Keeper) calculateBurnBySwapOut(
	ctx sdk.Context,
	burnIn sdk.Coin,
	backingDenom string,
) (
	backingOut sdk.Coin,
	ironOut sdk.Coin,
	burnFee sdk.Coin,
	err error,
) {
	err = k.checkBurnPriceUpperBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	ironPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.AttoIronDenom)
	if err != nil {
		return
	}

	backingRatio := k.GetBackingRatio(ctx)

	burnFee = computeFee(burnIn, backingParams.BurnFee)
	burnActual := burnIn.Sub(burnFee)
	burnActualInUSD := burnActual.Amount.ToDec().Mul(gridiron.MicroUSMTarget)

	backingOut = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	ironOut = sdk.NewCoin(gridiron.AttoIronDenom, sdk.ZeroInt())

	if backingRatio.GTE(sdk.OneDec()) {
		// full/over backing
		backingOut.Amount = burnActualInUSD.QuoTruncate(backingPrice).TruncateInt()
	} else if backingRatio.IsZero() {
		// full algorithmic
		ironOut.Amount = burnActualInUSD.QuoTruncate(ironPrice).TruncateInt()
	} else {
		// fractional
		backingOut.Amount = burnActualInUSD.Mul(backingRatio).QuoTruncate(backingPrice).TruncateInt()
		ironOut.Amount = burnActualInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoTruncate(ironPrice).TruncateInt()
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)

	poolBackingBalance := sdk.NewCoin(backingDenom, sdk.MinInt(poolBacking.Backing.Amount, moduleOwnedBacking.Amount))
	if poolBackingBalance.IsLT(backingOut) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) > balance(%s)", backingOut, poolBackingBalance)
		return
	}

	return
}

func (k Keeper) calculateBuyBackingIn(
	ctx sdk.Context,
	backingOut sdk.Coin,
) (
	ironIn sdk.Coin,
	buybackFee sdk.Coin,
	err error,
) {
	backingDenom := backingOut.Denom

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	ironPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.AttoIronDenom)
	if err != nil {
		return
	}

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}

	backingOutTotal := sdk.NewCoin(backingDenom, backingOut.Amount.ToDec().Quo(sdk.OneDec().Sub(*backingParams.BuybackFee)).TruncateInt())
	ironInValue := backingOutTotal.Amount.ToDec().Mul(backingPrice)

	if ironInValue.GT(excessBackingValue.ToDec()) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "")
		return
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)

	poolBackingBalance := sdk.NewCoin(backingDenom, sdk.MinInt(poolBacking.Backing.Amount, moduleOwnedBacking.Amount))
	if poolBackingBalance.IsLT(backingOutTotal) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) > balance(%s)", backingOutTotal, poolBackingBalance)
		return
	}

	ironIn = sdk.NewCoin(gridiron.AttoIronDenom, ironInValue.Quo(ironPrice).RoundInt())
	buybackFee = sdk.NewCoin(backingDenom, backingOutTotal.Amount.ToDec().Mul(*backingParams.BuybackFee).RoundInt())
	return
}

func (k Keeper) calculateBuyBackingOut(
	ctx sdk.Context,
	ironIn sdk.Coin,
	backingDenom string,
) (
	backingOut sdk.Coin,
	buybackFee sdk.Coin,
	err error,
) {
	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	ironPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.AttoIronDenom)
	if err != nil {
		return
	}

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}

	ironInValue := ironIn.Amount.ToDec().Mul(ironPrice)
	if ironInValue.GT(excessBackingValue.ToDec()) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "")
		return
	}

	backingOutTotal := sdk.NewCoin(backingDenom, ironInValue.Quo(backingPrice).TruncateInt())

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)

	poolBackingBalance := sdk.NewCoin(backingDenom, sdk.MinInt(poolBacking.Backing.Amount, moduleOwnedBacking.Amount))
	if poolBackingBalance.IsLT(backingOutTotal) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) > balance(%s)", backingOutTotal, poolBackingBalance)
		return
	}

	buybackFee = computeFee(backingOutTotal, backingParams.BuybackFee)
	backingOut = backingOutTotal.Sub(buybackFee)
	return
}

func (k Keeper) calculateSellBackingIn(
	ctx sdk.Context,
	ironOut sdk.Coin,
	backingDenom string,
) (
	backingIn sdk.Coin,
	rebackFee sdk.Coin,
	err error,
) {
	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	ironPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.AttoIronDenom)
	if err != nil {
		return
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}
	missingBackingValue := excessBackingValue.Neg()
	availableIronMint := missingBackingValue.ToDec().Quo(ironPrice)

	bonusRatio := k.RebackBonus(ctx)

	ironMint := ironOut.Amount.ToDec().Quo(sdk.OneDec().Add(bonusRatio).Sub(*backingParams.RebackFee))

	backingIn = sdk.NewCoin(backingDenom, ironMint.Mul(ironPrice).Quo(backingPrice).RoundInt())
	rebackFee = sdk.NewCoin(gridiron.AttoIronDenom, ironMint.Mul(*backingParams.RebackFee).RoundInt())

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "")
		return
	}
	if ironMint.GT(availableIronMint) {
		err = sdkerrors.Wrap(types.ErrIronCoinInsufficient, "")
		return
	}

	return
}

func (k Keeper) calculateSellBackingOut(
	ctx sdk.Context,
	backingIn sdk.Coin,
) (
	ironOut sdk.Coin,
	rebackFee sdk.Coin,
	err error,
) {
	backingDenom := backingIn.Denom

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	ironPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.AttoIronDenom)
	if err != nil {
		return
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "")
		return
	}

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}
	missingBackingValue := excessBackingValue.Neg()
	availableIronMint := missingBackingValue.ToDec().Quo(ironPrice)

	bonusRatio := k.RebackBonus(ctx)
	ironMint := sdk.NewCoin(gridiron.AttoIronDenom, backingIn.Amount.ToDec().Mul(backingPrice).Quo(ironPrice).TruncateInt())
	bonus := computeFee(ironMint, &bonusRatio)
	rebackFee = computeFee(ironMint, backingParams.RebackFee)

	if ironMint.Amount.ToDec().GT(availableIronMint) {
		err = sdkerrors.Wrap(types.ErrIronCoinInsufficient, "")
		return
	}

	ironOut = ironMint.Add(bonus).Sub(rebackFee)
	return
}

func (k Keeper) calculateMintByCollateral(
	ctx sdk.Context,
	account sdk.AccAddress,
	collateralDenom string,
	mintOut sdk.Coin,
) (
	mintFee sdk.Coin,
	totalColl types.TotalCollateral,
	poolColl types.PoolCollateral,
	accColl types.AccountCollateral,
	err error,
) {
	collateralParams, err := k.getAvailableCollateralParams(ctx, collateralDenom)
	if err != nil {
		return
	}

	// get prices in usd
	collateralPrice, err := k.oracleKeeper.GetExchangeRate(ctx, collateralDenom)
	if err != nil {
		return
	}
	ironPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.AttoIronDenom)
	if err != nil {
		return
	}

	totalColl, poolColl, accColl, err = k.getCollateral(ctx, account, collateralDenom)
	if err != nil {
		return
	}

	// settle interest fee
	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, *collateralParams.InterestFee)

	// compute mint total
	mintFee = computeFee(mintOut, collateralParams.MintFee)
	mintTotal := mintOut.Add(mintFee)

	// update grid debt
	accColl.GridDebt = accColl.GridDebt.Add(mintTotal)
	poolColl.GridDebt = poolColl.GridDebt.Add(mintTotal)
	totalColl.GridDebt = totalColl.GridDebt.Add(mintTotal)

	if collateralParams.MaxGridMint != nil && poolColl.GridDebt.Amount.GT(*collateralParams.MaxGridMint) {
		err = sdkerrors.Wrapf(types.ErrGridCeiling, "")
		return
	}

	collateralValue := accColl.Collateral.Amount.ToDec().Mul(collateralPrice)
	ironCollateralizedValue := accColl.IronCollateralized.Amount.ToDec().Mul(ironPrice)
	if !collateralValue.IsPositive() {
		err = sdkerrors.Wrapf(types.ErrAccountInsufficientCollateral, "")
		return
	}

	actualCatalyticRatio := sdk.MinDec(ironCollateralizedValue.Quo(collateralValue), *collateralParams.CatalyticIronRatio)

	// actualCatalyticRatio / catalyticRatio = (availableLTV - basicLTV) / (maxLTV - basicLTV)
	availableLTV := *collateralParams.BasicLoanToValue
	if collateralParams.CatalyticIronRatio.IsPositive() {
		availableLTV = availableLTV.Add(actualCatalyticRatio.Mul(collateralParams.LoanToValue.Sub(*collateralParams.BasicLoanToValue)).Quo(*collateralParams.CatalyticIronRatio))
	}
	availableDebtMax := collateralValue.Mul(availableLTV).Quo(gridiron.MicroUSMTarget).TruncateInt()

	if availableDebtMax.LT(accColl.GridDebt.Amount) {
		err = sdkerrors.Wrapf(types.ErrAccountInsufficientCollateral, "")
		return
	}

	return
}

func computeFee(coin sdk.Coin, rate *sdk.Dec) sdk.Coin {
	amt := sdk.ZeroInt()
	if rate != nil {
		amt = coin.Amount.ToDec().Mul(*rate).RoundInt()
	}
	return sdk.NewCoin(coin.Denom, amt)
}

func (k Keeper) checkMintPriceLowerBound(ctx sdk.Context) error {
	gridPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.MicroUSMDenom)
	if err != nil {
		return err
	}
	// market price must be >= target price + mint bias
	mintPriceLowerBound := gridiron.MicroUSMTarget.Mul(sdk.OneDec().Add(k.MintPriceBias(ctx)))
	if gridPrice.LT(mintPriceLowerBound) {
		return sdkerrors.Wrapf(types.ErrGridPriceTooLow, "%s price too low: %s", gridiron.MicroUSMDenom, gridPrice)
	}
	return nil
}

func (k Keeper) checkBurnPriceUpperBound(ctx sdk.Context) error {
	gridPrice, err := k.oracleKeeper.GetExchangeRate(ctx, gridiron.MicroUSMDenom)
	if err != nil {
		return err
	}
	// market price must be <= target price - burn bias
	burnPriceUpperBound := gridiron.MicroUSMTarget.Mul(sdk.OneDec().Sub(k.BurnPriceBias(ctx)))
	if gridPrice.GT(burnPriceUpperBound) {
		return sdkerrors.Wrapf(types.ErrGridPriceTooHigh, "%s price too high: %s", gridiron.MicroUSMDenom, gridPrice)
	}
	return nil
}

func (k Keeper) getAvailableBackingParams(ctx sdk.Context, backingDenom string) (backingParams types.BackingRiskParams, err error) {
	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}
	return
}

func (k Keeper) getAvailableCollateralParams(ctx sdk.Context, collateralDenom string) (collateralParams types.CollateralRiskParams, err error) {
	collateralParams, found := k.GetCollateralRiskParams(ctx, collateralDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", collateralDenom)
		return
	}
	if !collateralParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", collateralDenom)
		return
	}
	return
}

func (k Keeper) getExcessBackingValue(ctx sdk.Context) (excessBackingValue sdk.Int, err error) {
	totalBacking, found := k.GetTotalBacking(ctx)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "total backing not found")
		return
	}

	backingRatio := k.GetBackingRatio(ctx)
	requiredBackingValue := totalBacking.GridMinted.Amount.ToDec().Mul(backingRatio).Ceil().TruncateInt()
	if requiredBackingValue.IsNegative() {
		requiredBackingValue = sdk.ZeroInt()
	}

	totalBackingValue, err := k.totalBackingInUSD(ctx)
	if err != nil {
		return
	}

	// may be negative
	excessBackingValue = totalBackingValue.Sub(requiredBackingValue)
	return
}

func (k Keeper) totalBackingInUSD(ctx sdk.Context) (sdk.Int, error) {
	totalBackingValue := sdk.ZeroDec()
	for _, pool := range k.GetAllPoolBacking(ctx) {
		// get price in usd
		backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, pool.Backing.Denom)
		if err != nil {
			return sdk.Int{}, err
		}
		totalBackingValue = totalBackingValue.Add(pool.Backing.Amount.ToDec().Mul(backingPrice))
	}
	return totalBackingValue.TruncateInt(), nil
}
