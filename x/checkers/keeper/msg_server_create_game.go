package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dan-l/checkers/x/checkers/rules"
	"github.com/dan-l/checkers/x/checkers/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}

	// Get new index for game
	newIndex := strconv.FormatUint(systemInfo.NextId, 10)

	// Create game object
	newGame := rules.New()
	storedGame := types.StoredGame{
		Index:     newIndex,
		Board:     newGame.String(),
		Turn:      rules.PieceStrings[newGame.Turn],
		Black:     msg.Black,
		Red:       msg.Red,
		MoveCount: 0,
	}

	// Validate game
	err := storedGame.Validate()
	if err != nil {
		// Instead of panic, return err
		// Prevent panic attack on nodes, while attacker needs to pay gas fee to spam
		return nil, err
	}

	// Store the valid game
	k.Keeper.SetStoredGame(ctx, storedGame)

	// Update the id counter in SystemInfo
	systemInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.GameCreatedEventType,
			sdk.NewAttribute(types.GameCreatedEventCreator, msg.Creator),
			sdk.NewAttribute(types.GameCreatedEventGameIndex, storedGame.Index),
			sdk.NewAttribute(types.GameCreatedEventBlack, msg.Black),
			sdk.NewAttribute(types.GameCreatedEventRed, msg.Red),
		),
	)

	return &types.MsgCreateGameResponse{
		GameIndex: newIndex,
	}, nil
}
