package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/CosmWasm/wasmd/x/message/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagEvents = "events"

	eventFormat = "{eventType}.{eventAttribute}={value}"
)

// GetQueryCmd returns the transaction commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the message module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetMessagesCmd(cdc))

	return cmd
}

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
func GetMessagesCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages [address]",
		Short: "Query account messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			// accGetter := authtypes.NewQueryBalanceParams(cliCtx)

			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryBalanceParams(key))
			if err != nil {
				return err
			}
			msg, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QuerierRoute), bs)

			var message types.MsgTextsSend

			err = cliCtx.Codec.UnmarshalJSON(msg, &message)

			return cliCtx.PrintOutput(message)
		},
	}

	return flags.GetCommands(cmd)[0]
}
