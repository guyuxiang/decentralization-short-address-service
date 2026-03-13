package cli

import (
	"time"

	"github.com/spf13/cobra"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"sas/x/sas"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
)

func GetCmdBuySUrl(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy-sUrl [sUrl] [amount]",
		Short: "bid for existing sUrl or claim new sUrl",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := sas.NewMsgBuySUrl(args[0], coins, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdSetLUrl(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set_lUrl [sUrl] [lUrl]",
		Short: "set the value associated with a name that you own",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			msg := sas.NewMsgSetLUrl(args[0], args[1], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdSetSell(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set_sell [sUrl] [isSell]",
		Short: "set the sell flag that you own",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			sellArg, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			msg := sas.NewMsgSetSell(args[0], sellArg, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdSetPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set_price [sUrl] [amount]",
		Short: "set the price associated with a name that you own",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}
			msg := sas.NewMsgSetPrice(args[0], coins, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdRenew(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "renew [sUrl] [days]",
		Short: "renew a short URL by adding more time (in days)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			days, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			duration := time.Duration(days*24) * time.Hour
			msg := sas.NewMsgRenew(args[0], duration, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdBuyEscrow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy-escrow [sUrl] [amount]",
		Short: "buy a short URL with escrow (safer transaction)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := sas.NewMsgBuySUrlEscrow(args[0], coins, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdConfirmEscrow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "confirm-escrow [sUrl]",
		Short: "confirm an escrow transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			msg := sas.NewMsgConfirmEscrow(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdCancelEscrow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cancel-escrow [sUrl]",
		Short: "cancel an escrow transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			msg := sas.NewMsgCancelEscrow(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdBatchSetLUrl(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "batch-set-lurl [sUrl] [url1,url2,...]",
		Short: "batch set multiple long URLs for A/B testing",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			urls := strings.Split(args[1], ",")
			msg := sas.NewMsgBatchSetLUrl(args[0], urls, cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdAddBlackList(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-blacklist [url] [domain|url]",
		Short: "add a URL or domain to blacklist",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			isDomain := args[1] == "domain"
			msg := sas.NewMsgAddBlackList(args[0], isDomain, cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
