package sas

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "sas" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetLUrl:
			return handleMsgSetLUrl(ctx, keeper, msg)
		case MsgBuySUrl:
			return handleMsgBuySUrl(ctx, keeper, msg)
		case MsgSetSell:
			return handleMsgSetSell(ctx, keeper, msg)
		case MsgSetPrice:
			return handleMsgSetPrice(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized sas Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSetLUrl(ctx sdk.Context, keeper Keeper, msg MsgSetLUrl) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl) {
		return sdk.ErrUnauthorized("Incorrect Surl").Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}

	keeper.SetLUrl(ctx, msg.SUrl, msg.LUrl)
	return sdk.Result{}
}

func handleMsgBuySUrl(ctx sdk.Context, keeper Keeper, msg MsgBuySUrl) sdk.Result {
	if len(msg.SUrl) != 0 {
		if CheckSUrlExist(ctx, keeper, msg.SUrl) {
			return txUrl(ctx, keeper, msg)
		} else {
			_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
			if err != nil {
				return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
			}
			keeper.StoreLAdress(ctx, msg.SUrl, msg.Buyer, msg.Bid)
			GlobalBloomFilter.Set(msg.SUrl)
		}
	} else {
		newSUrl := ApplyShortUrl(ctx, keeper)
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			rullBackNumber()
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
		keeper.StoreLAdress(ctx, newSUrl, msg.Buyer, msg.Bid)
		GlobalBloomFilter.Set(newSUrl)
	}
	return sdk.Result{}
}

func txUrl(ctx sdk.Context, keeper Keeper, msg MsgBuySUrl) sdk.Result {
	if keeper.HasOwner(ctx, msg.SUrl) && !keeper.GetSell(ctx, msg.SUrl) {
		return sdk.ErrInternal("Adress dose not sell").Result()
	}
	if keeper.GetPrice(ctx, msg.SUrl).IsAllGT(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.SUrl), msg.Bid)
	if err != nil {
		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	}

	keeper.SetOwner(ctx, msg.SUrl, msg.Buyer)
	keeper.SetPrice(ctx, msg.SUrl, msg.Bid)
	return sdk.Result{}
}

func handleMsgSetSell(ctx sdk.Context, keeper Keeper, msg MsgSetSell) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl) {
		return sdk.ErrUnauthorized("Incorrect Surl").Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	keeper.SetSell(ctx, msg.SUrl, msg.IsSell)

	return sdk.Result{}
}

func handleMsgSetPrice(ctx sdk.Context, keeper Keeper, msg MsgSetPrice) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl) {
		return sdk.ErrUnauthorized("Incorrect Surl").Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	keeper.SetPrice(ctx, msg.SUrl, msg.Bid)
	return sdk.Result{}
}
