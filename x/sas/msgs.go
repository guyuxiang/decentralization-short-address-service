package sas

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSetLUrl struct {
	SUrl  string
	LUrl  string
	Owner sdk.AccAddress
}

func NewMsgSetLUrl(sUrl string, lUrl string, owner sdk.AccAddress) MsgSetLUrl {
	return MsgSetLUrl{
		SUrl:  sUrl,
		LUrl:  lUrl,
		Owner: owner,
	}
}

func (msg MsgSetLUrl) Route() string { return "sas" }

// Type should return the action
func (msg MsgSetLUrl) Type() string { return "set_lUrl" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetLUrl) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.SUrl) == 0 || len(msg.LUrl) == 0 {
		return sdk.ErrUnknownRequest("SUrl and/or LUrl cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetLUrl) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSetLUrl) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgBuySUrl struct {
	SUrl  string
	Bid   sdk.Coins
	Buyer sdk.AccAddress
}

func NewMsgBuySUrl(sUrl string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuySUrl {
	return MsgBuySUrl{
		SUrl:  sUrl,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Route should return the name of the module
func (msg MsgBuySUrl) Route() string { return "sas" }

// Type should return the action
func (msg MsgBuySUrl) Type() string { return "buy_sUrl" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBuySUrl) ValidateBasic() sdk.Error {
	if len(msg.SUrl) > 0 && len(msg.SUrl) != 6 {
		return sdk.ErrInvalidAddress(msg.SUrl)
	}
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBuySUrl) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgBuySUrl) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

type MsgSetSell struct {
	SUrl   string
	IsSell bool
	Owner  sdk.AccAddress
}

func NewMsgSetSell(sUrl string, isSell bool, owner sdk.AccAddress) MsgSetSell {
	return MsgSetSell{
		SUrl:   sUrl,
		IsSell: isSell,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgSetSell) Route() string { return "sas" }

// Type should return the action
func (msg MsgSetSell) Type() string { return "set_sell" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetSell) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.SUrl) == 0 {
		return sdk.ErrUnknownRequest("SUrl and/or LUrl cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetSell) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSetSell) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgSetPrice struct {
	SUrl  string
	Bid   sdk.Coins
	Owner sdk.AccAddress
}

func NewMsgSetPrice(sUrl string, bid sdk.Coins, owner sdk.AccAddress) MsgSetPrice {
	return MsgSetPrice{
		SUrl:  sUrl,
		Bid:   bid,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgSetPrice) Route() string { return "sas" }

// Type should return the action
func (msg MsgSetPrice) Type() string { return "set_price" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetPrice) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.SUrl) == 0 {
		return sdk.ErrUnknownRequest("SUrl and/or LUrl cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetPrice) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSetPrice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
