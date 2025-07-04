package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"time"

	"github.com/atomone-hub/atomone/x/gov/types"
)

var _, _ sdk.Msg = &MsgSendSecureText{}, &MsgUpdateParams{}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateParams) Route() string { return types.RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateParams) Type() string { return sdk.MsgTypeURL(&msg) }

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	return msg.Params.ValidateBasic()
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams.
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}

func NewMsgSendSecureText(encryptedFromAddr []byte, toAddr string, encryptedText []byte, fromAddr string) *MsgSendSecureText {
	now := time.Now()
	msg := SecureMessage{SubmitTime: &now,
		EncryptedFromAddress: encryptedFromAddr,
		ToAddress:            toAddr,
		EncryptedText:        encryptedText,
	}
	return &MsgSendSecureText{Msg: msg, From: fromAddr}
}

// Route implements the sdk.Msg interface.
func (msg MsgSendSecureText) Route() string { return types.RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgSendSecureText) Type() string { return sdk.MsgTypeURL(&msg) }

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgSendSecureText) ValidateBasic() error {
	return nil
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgSendSecureText) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams.
func (msg MsgSendSecureText) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.From)
	return []sdk.AccAddress{from}
}
