package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/atomone-hub/atomone/x/mymodule/types"
)

type Keeper struct {
	cdc       codec.BinaryCodec
	storeKey  storetypes.StoreKey
	authority string

	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
	stakingKeeper types.StakingKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authority string,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	stakingKeeper types.StakingKeeper,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		authority:     authority,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		stakingKeeper: stakingKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
