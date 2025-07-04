package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/atomone-hub/atomone/x/mymodule/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetMessages(ctx sdk.Context, to_address sdk.AccAddress) (msgs types.SecureMessages) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(to_address.Bytes())
	if bz == nil {
		return msgs
	}
	k.cdc.MustUnmarshal(bz, &msgs)
	return msgs
}

// SetParams set the params
func (k Keeper) AddMessage(ctx sdk.Context, to_address sdk.AccAddress, msg types.SecureMessage) error {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(to_address.Bytes())
	var msgs types.SecureMessages
	if bz == nil {
		msgs.Msgs = append(msgs.Msgs, &msg)
	} else {
		k.cdc.MustUnmarshal(bz, &msgs)
		msgs.Msgs = append(msgs.Msgs, &msg)
	}
	bz, err := k.cdc.Marshal(&msgs)
	if err != nil {
		return err
	}
	store.Set(to_address.Bytes(), bz)
	return nil
}
