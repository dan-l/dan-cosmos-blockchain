package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/checkers module sentinel errors
var (
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")
)

// x/checkers module storedGame errors
var (
	ErrInvalidBlack     = sdkerrors.Register(ModuleName, 1101, "black address is invalid: %s")
	ErrInvalidRed       = sdkerrors.Register(ModuleName, 1102, "red address is invalid: %s")
	ErrGameNotParseable = sdkerrors.Register(ModuleName, 1103, "game cannot be parsed")
)

// x/checkers module playMove errors
var (
	ErrGameNotFound     = sdkerrors.Register(ModuleName, 1104, "game id not found")
	ErrCreatorNotPlayer = sdkerrors.Register(ModuleName, 1105, "message creator is not a player")
	ErrNotPlayerTurn    = sdkerrors.Register(ModuleName, 1106, "player tried to play out of turn")
	ErrWrongMove        = sdkerrors.Register(ModuleName, 1107, "wrong move")
)
