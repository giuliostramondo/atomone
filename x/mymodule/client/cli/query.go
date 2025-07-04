package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/atomone-hub/atomone/x/mymodule/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group photon queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetQueryParamsCmd(),
		GetQueryMessagesCmd(),
	)
	return cmd
}

func GetQueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func Decrypt(clientCtx client.Context, cmd *cobra.Command, addrStr string, ciphertext []byte) (string, error) {
	addr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		err := fmt.Errorf("Not a valid AtomOne address: %s", addrStr)
		return "", err
	}
	key, err := clientCtx.Keyring.KeyByAddress(addr)
	if err != nil {
		err := fmt.Errorf("Could not find address in keyring: %s", addrStr)
		return "", err
	}
	priv, ok := key.GetLocal().PrivKey.GetCachedValue().(cryptotypes.PrivKey)
	if ok != true {
		err := fmt.Errorf("Could not access private key of: %s", addrStr)
		return "", err
	}
	privKey, _ := secp256k1.PrivKeyFromBytes(priv.Bytes())
	plaintext, err := secp256k1.Decrypt(privKey, ciphertext)
	return string(plaintext), err

}

func GetQueryMessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages [account]",
		Short: "Retrive account messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Messages(cmd.Context(),
				&types.QueryMessagesRequest{ToAddress: args[0]})
			if err != nil {
				return err
			}

			if len(res.Msgs.Msgs) == 0 {
				fmt.Println("No Messages.")
				return nil
			}
			for i := range res.Msgs.Msgs {
				encryptedFromAddress := res.Msgs.Msgs[i].EncryptedFromAddress
				fromAddress, err := Decrypt(clientCtx, cmd, args[0], encryptedFromAddress)
				if err != nil {
					return err
				}
				encryptedText := res.Msgs.Msgs[i].EncryptedText
				text, err := Decrypt(clientCtx, cmd, args[0], encryptedText)
				if err != nil {
					return err
				}
				fmt.Println("-------------")
				fmt.Println("Message: ", i)
				fmt.Println("From: ", fromAddress)
				fmt.Println("To: ", res.Msgs.Msgs[i].ToAddress)
				fmt.Println("Date: ", res.Msgs.Msgs[i].SubmitTime.Format(time.RFC822))
				fmt.Println("Message: ")
				fmt.Println(text)
			}
			return nil
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
