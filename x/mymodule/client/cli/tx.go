package cli

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/atomone-hub/atomone/x/mymodule/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/decred/dcrd/dcrec/secp256k1"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(GetTxSendSecureTextCmd())
	return cmd
}

// https://stackoverflow.com/questions/49475216/use-secp256k1-in-go
// https://github.com/decred/dcrd/tree/master/dcrec/secp256k1
func Encrypt(clientCtx client.Context, cmd *cobra.Command, addrStr string, plaintext string) ([]byte, error) {
	toAddr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		err := fmt.Errorf("Not a valid AtomOne address: %s", addrStr)
		return nil, err
	}
	queryClient := authtypes.NewQueryClient(clientCtx)
	res, err := queryClient.Account(cmd.Context(), &authtypes.QueryAccountRequest{Address: toAddr.String()})
	if err != nil {
		err := fmt.Errorf("Cannot retrive the public key of the account: %s", addrStr)
		return nil, err
	}
	var acc authtypes.AccountI
	clientCtx.Codec.UnpackAny(res.Account, &acc)
	pubKey := acc.GetPubKey()
	if pubKey == nil {
		err := fmt.Errorf("Cannot retrive the public key of the account: %s", addrStr)
		return nil, err
	}
	toAddr_pubKey_secp, err := secp256k1.ParsePubKey(acc.GetPubKey().Bytes())
	if err != nil {
		return nil, err
	}
	ciphertext, err := secp256k1.Encrypt(toAddr_pubKey_secp, []byte(plaintext))
	return ciphertext, err
}

func GetTxSendSecureTextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to_address] [message]",
		Short: "Send a secure text message to a cosmos address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			encryptedSenderAcc, err := Encrypt(clientCtx, cmd, args[0], clientCtx.GetFromAddress().String())
			if err != nil {
				return err
			}
			encryptedMessage, err := Encrypt(clientCtx, cmd, args[0], args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgSendSecureText(encryptedSenderAcc, args[0], encryptedMessage, clientCtx.GetFromAddress().String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
