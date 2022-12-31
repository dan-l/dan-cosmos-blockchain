package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dan-l/checkers/testutil/keeper"
	"github.com/dan-l/checkers/x/checkers"
	"github.com/dan-l/checkers/x/checkers/keeper"
	"github.com/dan-l/checkers/x/checkers/testutil"
	"github.com/dan-l/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

const (
	alice = testutil.Alice
	bob   = testutil.Bob
	carol = testutil.Carol
)

func TestCreateGame(t *testing.T) {
	msgServer, keeper, ctx := setupMsgServerCreateGame(t)
	createResponse, err := msgServer.CreateGame(
		ctx,
		&types.MsgCreateGame{
			Creator: alice,
			Black:   bob,
			Red:     carol,
		},
	)
	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "1",
	}, *createResponse)

	// Check state
	systemInfo, _ := keeper.GetSystemInfo(sdk.UnwrapSDKContext(ctx))
	require.EqualValues(t, types.SystemInfo{
		NextId: 2,
	}, systemInfo)
	game, _ := keeper.GetStoredGame(sdk.UnwrapSDKContext(ctx), createResponse.GameIndex)
	require.NotNil(t, game)
	require.Equal(t, bob, game.Black)
	require.Equal(t, carol, game.Red)
	require.Equal(t, uint64(0), game.MoveCount)

	// Check emit
	unwrapped_ctx := sdk.UnwrapSDKContext(ctx)
	events := sdk.StringifyEvents(unwrapped_ctx.EventManager().ABCIEvents())
	event := events[len(events)-1]
	require.EqualValues(
		t,
		sdk.StringEvent{
			Type: types.GameCreatedEventType,
			Attributes: []sdk.Attribute{
				{Key: "creator", Value: alice},
				{Key: "game-index", Value: createResponse.GameIndex},
				{Key: "black", Value: bob},
				{Key: "red", Value: carol},
			},
		},
		event,
	)
}

func TestCreate3Game(t *testing.T) {
	msgServer, keeper, ctx := setupMsgServerCreateGame(t)
	// Game 1
	createResponse, err := msgServer.CreateGame(
		ctx,
		&types.MsgCreateGame{
			Creator: alice,
			Black:   bob,
			Red:     carol,
		},
	)
	// Check Game 1
	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "1",
	}, *createResponse)
	game, _ := keeper.GetStoredGame(sdk.UnwrapSDKContext(ctx), createResponse.GameIndex)
	require.NotNil(t, game)
	require.Equal(t, bob, game.Black)
	require.Equal(t, carol, game.Red)

	// Game 2
	createResponse2, err2 := msgServer.CreateGame(
		ctx,
		&types.MsgCreateGame{
			Creator: bob,
			Black:   alice,
			Red:     carol,
		},
	)
	// Check Game 2
	require.Nil(t, err2)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "2",
	}, *createResponse2)
	game2, _ := keeper.GetStoredGame(sdk.UnwrapSDKContext(ctx), createResponse2.GameIndex)
	require.NotNil(t, game2)
	require.Equal(t, alice, game2.Black)
	require.Equal(t, carol, game2.Red)

	// Game 3
	createResponse3, err3 := msgServer.CreateGame(
		ctx,
		&types.MsgCreateGame{
			Creator: carol,
			Black:   alice,
			Red:     bob,
		},
	)
	// Check Game 3
	require.Nil(t, err3)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "3",
	}, *createResponse3)
	game3, _ := keeper.GetStoredGame(sdk.UnwrapSDKContext(ctx), createResponse3.GameIndex)
	require.NotNil(t, game3)
	require.Equal(t, alice, game3.Black)
	require.Equal(t, bob, game3.Red)

	games := keeper.GetAllStoredGame(sdk.UnwrapSDKContext(ctx))
	require.Len(t, games, 3)
}

func TestCreateGameBadRedAddress(t *testing.T) {
	msgServer, keeper, ctx := setupMsgServerCreateGame(t)
	createResponse, err := msgServer.CreateGame(
		ctx,
		&types.MsgCreateGame{
			Creator: alice,
			Black:   bob,
			Red:     "0x00",
		},
	)
	require.Nil(t, createResponse)
	require.ErrorContains(t, err, "red address is invalid")
	// No games created
	games := keeper.GetAllStoredGame(sdk.UnwrapSDKContext(ctx))
	require.Len(t, games, 0)
}

func TestCreateGameMissingRedAddress(t *testing.T) {
	msgServer, keeper, ctx := setupMsgServerCreateGame(t)
	createResponse, err := msgServer.CreateGame(
		ctx,
		&types.MsgCreateGame{
			Creator: alice,
			Black:   bob,
		},
	)
	require.Nil(t, createResponse)
	require.ErrorContains(t, err, "red address is invalid")
	// No games created
	games := keeper.GetAllStoredGame(sdk.UnwrapSDKContext(ctx))
	require.Len(t, games, 0)
}

func setupMsgServerCreateGame(t *testing.T) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	// Init keeper with genesis
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	return keeper.NewMsgServerImpl(*k), *k, sdk.WrapSDKContext(ctx)
}
