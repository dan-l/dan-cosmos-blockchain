package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dan-l/checkers/x/checkers/rules"
	"github.com/dan-l/checkers/x/checkers/types"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
	}

	// Validate player
	isBlack := storedGame.Black == msg.Creator
	isRed := storedGame.Red == msg.Creator
	if !isBlack && !isRed {
		return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
	}

	// Validate player turn
	var player rules.Player
	if isBlack {
		player = rules.BLACK_PLAYER
	} else {
		player = rules.RED_PLAYER
	}
	game, err := storedGame.ParseGame()
	if err != nil {
		// DB might be corrupted, panic
		panic(err.Error())
	}
	if !game.TurnIs(player) {
		return nil, sdkerrors.Wrapf(types.ErrNotPlayerTurn, "%s", player)
	}

	// Move
	captured, moveErr := game.Move(
		rules.Pos{
			X: int(msg.FromX),
			Y: int(msg.FromY),
		},
		rules.Pos{
			X: int(msg.ToX),
			Y: int(msg.ToY),
		},
	)
	if moveErr != nil {
		return nil, sdkerrors.Wrapf(types.ErrWrongMove, moveErr.Error())
	}

	// Update board
	storedGame.MoveCount++
	storedGame.Board = game.String()
	storedGame.Turn = rules.PieceStrings[game.Turn]
	k.Keeper.SetStoredGame(ctx, storedGame)

	winner := rules.PieceStrings[game.Winner()]
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.PlayMovedEventType,
			sdk.NewAttribute(types.PlayMovedEventCreator, msg.Creator),
			sdk.NewAttribute(types.PlayMovedEventGameIndex, storedGame.Index),
			sdk.NewAttribute(types.PlayMovedEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
			sdk.NewAttribute(types.PlayMovedEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
			sdk.NewAttribute(types.PlayMovedEventWinner, winner),
		),
	)

	return &types.MsgPlayMoveResponse{
		CapturedX: int32(captured.X),
		CapturedY: int32(captured.Y),
		Winner:    winner,
	}, nil
}
