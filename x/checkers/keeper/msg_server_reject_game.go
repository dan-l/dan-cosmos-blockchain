package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dan-l/checkers/x/checkers/types"
)

func (k msgServer) RejectGame(goCtx context.Context, msg *types.MsgRejectGame) (*types.MsgRejectGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
	}

	isBlack := storedGame.Black == msg.Creator
	isRed := storedGame.Red == msg.Creator
	if !isBlack && !isRed {
		return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
	}

	// Check if moved
	if isBlack && storedGame.MoveCount > 0 {
		return nil, types.ErrBlackAlreadyPlayed
	}
	if isRed && storedGame.MoveCount > 1 {
		return nil, types.ErrRedAlreadyPlayed
	}

	k.Keeper.RemoveStoredGame(ctx, msg.GameIndex)

	// Emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.RejectGameEventType,
			sdk.NewAttribute(types.RejectGameEventCreator, msg.Creator),
			sdk.NewAttribute(types.RejectGameEventGameIndex, msg.GameIndex),
		),
	)

	return &types.MsgRejectGameResponse{}, nil
}
