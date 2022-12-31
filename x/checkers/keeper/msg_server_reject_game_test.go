package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dan-l/checkers/testutil/keeper"
	"github.com/dan-l/checkers/x/checkers"
	"github.com/dan-l/checkers/x/checkers/keeper"
	"github.com/dan-l/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

func TestRejectGame(t *testing.T) {
	msgServer, keeper, ctx, gameIndex := setupMsgServerForRejectGame(t)
	_, err := msgServer.RejectGame(
		ctx,
		&types.MsgRejectGame{
			Creator:   bob,
			GameIndex: gameIndex,
		},
	)

	// Check
	require.Nil(t, err)
	unwrapped_ctx := sdk.UnwrapSDKContext(ctx)
	systemInfo, _ := keeper.GetSystemInfo(unwrapped_ctx)
	require.EqualValues(t, 2, systemInfo.NextId)
	_, found := keeper.GetStoredGame(unwrapped_ctx, gameIndex)
	require.False(t, found)

	// Events
	events := sdk.StringifyEvents(unwrapped_ctx.EventManager().ABCIEvents())
	event := events[len(events)-1]
	require.EqualValues(
		t,
		sdk.StringEvent{
			Type: types.RejectGameEventType,
			Attributes: []sdk.Attribute{
				{Key: "creator", Value: bob},
				{Key: "game-index", Value: gameIndex},
			},
		},
		event,
	)
}

func TestRejectGameBlackCannotReject(t *testing.T) {
	msgServer, keeper, ctx, gameIndex := setupMsgServerForRejectGame(t)

	// Move
	msgServer.PlayMove(
		ctx,
		&types.MsgPlayMove{
			Creator:   bob,
			GameIndex: gameIndex,
			FromX:     1,
			FromY:     2,
			ToX:       2,
			ToY:       3,
		},
	)

	// Reject
	rejectGameResponse, err := msgServer.RejectGame(
		ctx,
		&types.MsgRejectGame{
			Creator:   bob,
			GameIndex: gameIndex,
		},
	)
	require.Nil(t, rejectGameResponse)
	require.ErrorContains(t, err, types.ErrBlackAlreadyPlayed.Error())
	unwrapped_ctx := sdk.UnwrapSDKContext(ctx)
	_, found := keeper.GetStoredGame(unwrapped_ctx, gameIndex)
	require.True(t, found)
}

func TestRejectGameRedCanReject(t *testing.T) {
	msgServer, keeper, ctx, gameIndex := setupMsgServerForRejectGame(t)

	// Move
	msgServer.PlayMove(
		ctx,
		&types.MsgPlayMove{
			Creator:   bob,
			GameIndex: gameIndex,
			FromX:     1,
			FromY:     2,
			ToX:       2,
			ToY:       3,
		},
	)

	// Reject
	_, err := msgServer.RejectGame(
		ctx,
		&types.MsgRejectGame{
			Creator:   carol,
			GameIndex: gameIndex,
		},
	)

	// Check
	require.Nil(t, err)
	unwrapped_ctx := sdk.UnwrapSDKContext(ctx)
	systemInfo, _ := keeper.GetSystemInfo(unwrapped_ctx)
	require.EqualValues(t, 2, systemInfo.NextId)
	_, found := keeper.GetStoredGame(unwrapped_ctx, gameIndex)
	require.False(t, found)
}

func TestRejectGameRedCannotReject(t *testing.T) {
	msgServer, keeper, ctx, gameIndex := setupMsgServerForRejectGame(t)

	// Move
	msgServer.PlayMove(
		ctx,
		&types.MsgPlayMove{
			Creator:   bob,
			GameIndex: gameIndex,
			FromX:     1,
			FromY:     2,
			ToX:       2,
			ToY:       3,
		},
	)
	msgServer.PlayMove(
		ctx,
		&types.MsgPlayMove{
			Creator:   carol,
			GameIndex: gameIndex,
			FromX:     0,
			FromY:     5,
			ToX:       1,
			ToY:       4,
		},
	)

	// Reject
	rejectGameResponse, err := msgServer.RejectGame(
		ctx,
		&types.MsgRejectGame{
			Creator:   carol,
			GameIndex: gameIndex,
		},
	)
	require.Nil(t, rejectGameResponse)
	require.ErrorContains(t, err, types.ErrRedAlreadyPlayed.Error())
	unwrapped_ctx := sdk.UnwrapSDKContext(ctx)
	_, found := keeper.GetStoredGame(unwrapped_ctx, gameIndex)
	require.True(t, found)
}

func TestRejectGameCreatorNotPlayer(t *testing.T) {
	msgServer, _, ctx, gameIndex := setupMsgServerForRejectGame(t)
	rejectGameResponse, err := msgServer.RejectGame(
		ctx,
		&types.MsgRejectGame{
			Creator:   alice,
			GameIndex: gameIndex,
		},
	)
	require.Nil(t, rejectGameResponse)
	require.ErrorContains(t, err, types.ErrCreatorNotPlayer.Error())
}

func setupMsgServerForRejectGame(t *testing.T) (types.MsgServer, keeper.Keeper, context.Context, string) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)
	game, _ := server.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	return server, *k, context, game.GameIndex
}
