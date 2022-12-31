package types_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dan-l/checkers/x/checkers/rules"
	"github.com/dan-l/checkers/x/checkers/testutil"
	"github.com/dan-l/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

const (
	alice = testutil.Alice
	bob   = testutil.Bob
)

func GetStoredGame1() types.StoredGame {
	return types.StoredGame{
		Black: alice,
		Red:   bob,
		Index: "1",
		Board: rules.New().String(),
		Turn:  "b",
	}
}

func TestCanGetAddressBlack(t *testing.T) {
	aliceAddress, err1 := sdk.AccAddressFromBech32(alice)
	black, err2 := GetStoredGame1().GetBlackAddress()
	require.Equal(t, aliceAddress, black)
	require.Nil(t, err1)
	require.Nil(t, err2)
}

func TestGetAddressWrongBlack(t *testing.T) {
	storedGame := GetStoredGame1()
	// Not a valid address
	storedGame.Black = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4"
	black, err := storedGame.GetBlackAddress()
	require.Nil(t, black)
	require.EqualError(
		t,
		storedGame.Validate(),
		err.Error(),
	)
}

func TestCanParseGame(t *testing.T) {
	game, err := GetStoredGame1().ParseGame()
	require.Nil(t, err)
	require.EqualValues(
		t,
		rules.New().String(),
		game.String(),
	)
}

func TestCanParseGameMove(t *testing.T) {
	storedGame := GetStoredGame1()
	// Update board move
	storedGame.Board = strings.Replace(storedGame.Board, "b", "r", 1)
	game, err := storedGame.ParseGame()
	require.Nil(t, err)
	require.NotEqualValues(
		t,
		rules.New().String(),
		game.String(),
	)
}

func TestCannotParseInvalidPlayer(t *testing.T) {
	storedGame := GetStoredGame1()
	// a is not a valid player
	storedGame.Board = strings.Replace(storedGame.Board, "b", "a", 1)
	game, err := storedGame.ParseGame()
	require.Nil(t, game)
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestCannotParseInvalidTurn(t *testing.T) {
	storedGame := GetStoredGame1()
	// a is not a valid player
	storedGame.Turn = "a"
	game, err := storedGame.ParseGame()
	require.Nil(t, game)
	require.EqualError(t, storedGame.Validate(), err.Error())
}
