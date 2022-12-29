package keeper

import (
	"github.com/dan-l/checkers/x/checkers/types"
)

var _ types.QueryServer = Keeper{}
