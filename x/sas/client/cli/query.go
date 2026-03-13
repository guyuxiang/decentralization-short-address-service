package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"sas/x/sas"
)

func GetCmdLUrl(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "lurl [sUrl]",
		Short: "Query long URL by short URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sUrl := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/LUrl/%s", queryRoute, sUrl), nil)
			if err != nil {
				fmt.Printf("could not resolve sUrl - %s \n", string(sUrl))
				return nil
			}

			var out sas.QueryResLUrl
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdLAddress(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "laddress [sUrl]",
		Short: "Query LAddress info of sUrl",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sUrl := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/LAddress/%s", queryRoute, sUrl), nil)
			if err != nil {
				fmt.Printf("could not resolve address - %s \n", string(sUrl))
				return nil
			}

			var out sas.LAddress
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdSUrls(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "surls",
		Short: "Query all short URLs with pagination",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/SUrls", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get query sUrls\n")
				return nil
			}

			var out sas.QueryResPage
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdOwnerSUrls(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "owner-surls [owner]",
		Short: "Query short URLs by owner address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			owner := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/owner/%s", queryRoute, owner), nil)
			if err != nil {
				fmt.Printf("could not get query owner sUrls\n")
				return nil
			}

			var out sas.QueryResPage
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
