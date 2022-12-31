package keeper_test

import (
	"context"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dan-l/checkers/testutil/keeper"
	"github.com/dan-l/checkers/x/checkers"
	"github.com/dan-l/checkers/x/checkers/keeper"
	"github.com/dan-l/checkers/x/checkers/rules"
	"github.com/dan-l/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

func TestPlayMove(t *testing.T) {
	msgServer, keeper, ctx, gameIndex := setupMsgServerPlayMove(t)

	playMoveResponse, err := msgServer.PlayMove(
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

	// Check state
	require.Nil(t, err)
	require.EqualValues(
		t,
		types.MsgPlayMoveResponse{
			CapturedX: -1,
			CapturedY: -1,
			Winner:    "*",
		},
		*playMoveResponse,
	)
	game1, found := keeper.GetStoredGame(sdk.UnwrapSDKContext(ctx), "1")
	require.True(t, found)
	require.EqualValues(
		t,
		"*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		game1.Board,
	)
	require.EqualValues(
		t,
		rules.PieceStrings[rules.RED_PLAYER],
		game1.Turn,
	)

	// Check events
	unwrapped_ctx := sdk.UnwrapSDKContext(ctx)
	events := sdk.StringifyEvents(unwrapped_ctx.EventManager().ABCIEvents())
	event := events[len(events)-1]
	require.EqualValues(
		t,
		sdk.StringEvent{
			Type: types.PlayMovedEventType,
			Attributes: []sdk.Attribute{
				{Key: "creator", Value: bob},
				{Key: "game-index", Value: gameIndex},
				{Key: "captured-x", Value: strconv.FormatInt(int64(playMoveResponse.CapturedX), 10)},
				{Key: "captured-y", Value: strconv.FormatInt(int64(playMoveResponse.CapturedY), 10)},
				{Key: "winner", Value: playMoveResponse.Winner},
			},
		},
		event,
	)
}

func TestGameBoardCorrupted(t *testing.T) {
	msgServer, keeper, ctx, gameIndex := setupMsgServerPlayMove(t)
	// Corrupt board
	storedGame, _ := keeper.GetStoredGame(sdk.UnwrapSDKContext(ctx), gameIndex)
	storedGame.Board = "not a board"
	keeper.SetStoredGame(sdk.UnwrapSDKContext(ctx), storedGame)
	// setup try catch panic
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic")
		require.Equal(t, r, "game cannot be parsed: invalid board string: not a board")
	}()
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
}

func TestPlayMoveGameNotFound(t *testing.T) {
	msgServer, _, ctx, _ := setupMsgServerPlayMove(t)

	playMoveResponse, err := msgServer.PlayMove(
		ctx,
		&types.MsgPlayMove{
			Creator:   bob,
			GameIndex: "0",
			FromX:     1,
			FromY:     2,
			ToX:       2,
			ToY:       3,
		},
	)
	require.Nil(t, playMoveResponse)
	require.ErrorContains(t, err, types.ErrGameNotFound.Error())
}

func TestPlayMoveCreatorNotPlayer(t *testing.T) {
	msgServer, _, ctx, gameIndex := setupMsgServerPlayMove(t)

	playMoveResponse, err := msgServer.PlayMove(
		ctx,
		&types.MsgPlayMove{
			Creator:   alice,
			GameIndex: gameIndex,
			FromX:     1,
			FromY:     2,
			ToX:       2,
			ToY:       3,
		},
	)
	require.Nil(t, playMoveResponse)
	require.ErrorContains(t, err, types.ErrCreatorNotPlayer.Error())
}

func TestPlayMoveNotPlayerTurn(t *testing.T) {
	msgServer, _, ctx, gameIndex := setupMsgServerPlayMove(t)

	playMoveResponse, err := msgServer.PlayMove(
		ctx,
		&types.MsgPlayMove{
			Creator:   carol,
			GameIndex: gameIndex,
			FromX:     1,
			FromY:     2,
			ToX:       2,
			ToY:       3,
		},
	)
	require.Nil(t, playMoveResponse)
	require.ErrorContains(t, err, types.ErrNotPlayerTurn.Error())
}

func TestPlayMoveWrongMove(t *testing.T) {
	msgServer, _, ctx, gameIndex := setupMsgServerPlayMove(t)

	playMoveResponse, err := msgServer.PlayMove(
		ctx,
		&types.MsgPlayMove{
			Creator:   bob,
			GameIndex: gameIndex,
			FromX:     1,
			FromY:     0,
			ToX:       0,
			ToY:       1,
		},
	)
	require.Nil(t, playMoveResponse)
	require.ErrorContains(t, err, types.ErrWrongMove.Error())
}

func setupMsgServerPlayMove(t *testing.T) (types.MsgServer, keeper.Keeper, context.Context, string) {
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
