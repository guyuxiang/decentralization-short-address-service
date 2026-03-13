package sas

import (
	"encoding/json"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MaxLUrlLength = 2048
)

func validateLUrl(lUrl string) sdk.Error {
	if len(lUrl) == 0 {
		return sdk.ErrUnknownRequest("LUrl cannot be empty")
	}
	if len(lUrl) > MaxLUrlLength {
		return sdk.ErrUnknownRequest("LUrl too long")
	}
	lUrl = strings.ToLower(lUrl)
	if !strings.HasPrefix(lUrl, "http://") && !strings.HasPrefix(lUrl, "https://") {
		return sdk.ErrUnknownRequest("LUrl must start with http:// or https://")
	}
	return nil
}

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
	if err := validateLUrl(msg.LUrl); err != nil {
		return err
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
	SUrl   string
	Bid    sdk.Coins
	Buyer  sdk.AccAddress
	Length int
}

func NewMsgBuySUrl(sUrl string, bid sdk.Coins, buyer sdk.AccAddress, length int) MsgBuySUrl {
	return MsgBuySUrl{
		SUrl:   sUrl,
		Bid:    bid,
		Buyer:  buyer,
		Length: length,
	}
}

// Route should return the name of the module
func (msg MsgBuySUrl) Route() string { return "sas" }

// Type should return the action
func (msg MsgBuySUrl) Type() string { return "buy_sUrl" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBuySUrl) ValidateBasic() sdk.Error {
	if len(msg.SUrl) > 0 && (len(msg.SUrl) < 1 || len(msg.SUrl) > 6) {
		return sdk.ErrInvalidAddress("SUrl length must be 1-6")
	}
	if msg.Length > 0 && (msg.Length < 1 || msg.Length > 6) {
		return sdk.ErrUnknownRequest("Length must be 1-6")
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

type MsgRenew struct {
	SUrl     string
	Duration time.Duration
	Owner    sdk.AccAddress
}

func NewMsgRenew(sUrl string, duration time.Duration, owner sdk.AccAddress) MsgRenew {
	return MsgRenew{
		SUrl:     sUrl,
		Duration: duration,
		Owner:    owner,
	}
}

func (msg MsgRenew) Route() string { return "sas" }

func (msg MsgRenew) Type() string { return "renew" }

func (msg MsgRenew) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.SUrl) == 0 {
		return sdk.ErrUnknownRequest("SUrl cannot be empty")
	}
	if msg.Duration <= 0 {
		return sdk.ErrUnknownRequest("Duration must be positive")
	}
	return nil
}

func (msg MsgRenew) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgRenew) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgBuySUrlEscrow struct {
	SUrl   string
	Amount sdk.Coins
	Buyer  sdk.AccAddress
}

func NewMsgBuySUrlEscrow(sUrl string, amount sdk.Coins, buyer sdk.AccAddress) MsgBuySUrlEscrow {
	return MsgBuySUrlEscrow{
		SUrl:   sUrl,
		Amount: amount,
		Buyer:  buyer,
	}
}

func (msg MsgBuySUrlEscrow) Route() string { return "sas" }

func (msg MsgBuySUrlEscrow) Type() string { return "buy_sUrl_escrow" }

func (msg MsgBuySUrlEscrow) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.SUrl) == 0 {
		return sdk.ErrUnknownRequest("SUrl cannot be empty")
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive")
	}
	return nil
}

func (msg MsgBuySUrlEscrow) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgBuySUrlEscrow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

type MsgConfirmEscrow struct {
	SUrl      string
	Confirmor sdk.AccAddress
}

func NewMsgConfirmEscrow(sUrl string, confirmor sdk.AccAddress) MsgConfirmEscrow {
	return MsgConfirmEscrow{
		SUrl:      sUrl,
		Confirmor: confirmor,
	}
}

func (msg MsgConfirmEscrow) Route() string { return "sas" }

func (msg MsgConfirmEscrow) Type() string { return "confirm_escrow" }

func (msg MsgConfirmEscrow) ValidateBasic() sdk.Error {
	if msg.Confirmor.Empty() {
		return sdk.ErrInvalidAddress(msg.Confirmor.String())
	}
	if len(msg.SUrl) == 0 {
		return sdk.ErrUnknownRequest("SUrl cannot be empty")
	}
	return nil
}

func (msg MsgConfirmEscrow) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgConfirmEscrow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Confirmor}
}

type MsgCancelEscrow struct {
	SUrl     string
	Canceler sdk.AccAddress
}

func NewMsgCancelEscrow(sUrl string, canceler sdk.AccAddress) MsgCancelEscrow {
	return MsgCancelEscrow{
		SUrl:     sUrl,
		Canceler: canceler,
	}
}

func (msg MsgCancelEscrow) Route() string { return "sas" }

func (msg MsgCancelEscrow) Type() string { return "cancel_escrow" }

func (msg MsgCancelEscrow) ValidateBasic() sdk.Error {
	if msg.Canceler.Empty() {
		return sdk.ErrInvalidAddress(msg.Canceler.String())
	}
	if len(msg.SUrl) == 0 {
		return sdk.ErrUnknownRequest("SUrl cannot be empty")
	}
	return nil
}

func (msg MsgCancelEscrow) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCancelEscrow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Canceler}
}

type MsgBatchSetLUrl struct {
	SUrl  string
	LUrls []string
	Owner sdk.AccAddress
}

func NewMsgBatchSetLUrl(sUrl string, lUrls []string, owner sdk.AccAddress) MsgBatchSetLUrl {
	return MsgBatchSetLUrl{
		SUrl:  sUrl,
		LUrls: lUrls,
		Owner: owner,
	}
}

func (msg MsgBatchSetLUrl) Route() string { return "sas" }

func (msg MsgBatchSetLUrl) Type() string { return "batch_set_lUrl" }

func (msg MsgBatchSetLUrl) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.SUrl) == 0 {
		return sdk.ErrUnknownRequest("SUrl cannot be empty")
	}
	if len(msg.LUrls) == 0 {
		return sdk.ErrUnknownRequest("At least one LUrl required")
	}
	if len(msg.LUrls) > 10 {
		return sdk.ErrUnknownRequest("Maximum 10 URLs per batch")
	}
	for _, lUrl := range msg.LUrls {
		if err := validateLUrl(lUrl); err != nil {
			return err
		}
	}
	return nil
}

func (msg MsgBatchSetLUrl) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgBatchSetLUrl) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgAddBlackList struct {
	URL      string
	IsDomain bool
	Admin    sdk.AccAddress
}

func NewMsgAddBlackList(url string, isDomain bool, admin sdk.AccAddress) MsgAddBlackList {
	return MsgAddBlackList{
		URL:      url,
		IsDomain: isDomain,
		Admin:    admin,
	}
}

func (msg MsgAddBlackList) Route() string { return "sas" }

func (msg MsgAddBlackList) Type() string { return "add_blacklist" }

func (msg MsgAddBlackList) ValidateBasic() sdk.Error {
	if msg.Admin.Empty() {
		return sdk.ErrInvalidAddress(msg.Admin.String())
	}
	if len(msg.URL) == 0 {
		return sdk.ErrUnknownRequest("URL cannot be empty")
	}
	return nil
}

func (msg MsgAddBlackList) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgAddBlackList) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Admin}
}
