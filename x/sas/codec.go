package sas

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetLUrl{}, "sas/SetLUrl", nil)
	cdc.RegisterConcrete(MsgSetSell{}, "sas/SetSell", nil)
	cdc.RegisterConcrete(MsgBuySUrl{}, "sas/BuySUrl", nil)
	cdc.RegisterConcrete(MsgSetPrice{}, "sas/SetPrice", nil)
	cdc.RegisterConcrete(MsgRenew{}, "sas/Renew", nil)
	cdc.RegisterConcrete(MsgBuySUrlEscrow{}, "sas/BuySUrlEscrow", nil)
	cdc.RegisterConcrete(MsgConfirmEscrow{}, "sas/ConfirmEscrow", nil)
	cdc.RegisterConcrete(MsgCancelEscrow{}, "sas/CancelEscrow", nil)
	cdc.RegisterConcrete(MsgBatchSetLUrl{}, "sas/BatchSetLUrl", nil)
	cdc.RegisterConcrete(MsgAddBlackList{}, "sas/AddBlackList", nil)
}
