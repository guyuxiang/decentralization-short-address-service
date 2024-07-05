package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	sascmd "sas/x/sas/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group sas queries under a subcommand
	namesvcQueryCmd := &cobra.Command{
		Use:   "sas",
		Short: "Querying commands for the sas module",
	}

	namesvcQueryCmd.AddCommand(client.GetCommands(
		sascmd.GetCmdLUrl(mc.storeKey, mc.cdc),
		sascmd.GetCmdLAdress(mc.storeKey, mc.cdc),
		sascmd.GetCmdSUrls(mc.storeKey, mc.cdc),
	)...)

	return namesvcQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	namesvcTxCmd := &cobra.Command{
		Use:   "sas",
		Short: "sas transactions subcommands",
	}

	namesvcTxCmd.AddCommand(client.PostCommands(
		sascmd.GetCmdBuySUrl(mc.cdc),
		sascmd.GetCmdSetLUrl(mc.cdc),
		sascmd.GetCmdSetSell(mc.cdc),
		sascmd.GetCmdSetPrice(mc.cdc),
	)...)

	return namesvcTxCmd
}
