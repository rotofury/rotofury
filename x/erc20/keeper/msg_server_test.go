package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/gridiron-zone/gridiron/testutil/keeper"
	"github.com/gridiron-zone/gridiron/x/erc20/keeper"
	"github.com/gridiron-zone/gridiron/x/erc20/types"
)

func setupMsgServer(t *testing.T) (types.MsgServer, context.Context) {
	k, ctx := keepertest.Erc20Keeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
