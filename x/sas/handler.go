package sas

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeBuySUrl       = "buy_sUrl"
	EventTypeSetLUrl       = "set_lUrl"
	EventTypeSetSell       = "set_sell"
	EventTypeSetPrice      = "set_price"
	EventTypeRenew         = "renew"
	EventTypeExpire        = "expire"

	AttributeKeySUrl       = "sUrl"
	AttributeKeyLUrl      = "lUrl"
	AttributeKeyOwner     = "owner"
	AttributeKeyBuyer     = "buyer"
	AttributeKeyPrice     = "price"
	AttributeKeyIsSell   = "isSell"
	AttributeKeyExpiration = "expiration"
)

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
		case MsgRenew:
			return handleMsgRenew(ctx, keeper, msg)
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
	if keeper.IsExpired(ctx, msg.SUrl) {
		return sdk.ErrUnauthorized("SUrl has expired").Result()
	}

	keeper.SetLUrl(ctx, msg.SUrl, msg.LUrl)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeSetLUrl,
			sdk.NewAttribute(AttributeKeySUrl, msg.SUrl),
			sdk.NewAttribute(AttributeKeyLUrl, msg.LUrl),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
		),
	)

	return sdk.Result{}
}

func handleMsgBuySUrl(ctx sdk.Context, keeper Keeper, msg MsgBuySUrl) sdk.Result {
	fee := calculateFee(msg.Bid)
	totalCost := msg.Bid.Add(fee)

	if len(msg.SUrl) != 0 {
		if CheckSUrlExist(ctx, keeper, msg.SUrl) {
			if keeper.IsExpired(ctx, msg.SUrl) {
				keeper.DeleteLAddress(ctx, msg.SUrl)
			} else {
				return txUrl(ctx, keeper, msg, fee)
			}
		}
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, totalCost)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
		keeper.StoreLAddress(ctx, msg.SUrl, msg.Buyer, msg.Bid, DefaultRentDuration)
		GlobalBloomFilter.Set(msg.SUrl)
	} else {
		newSUrl := ApplyShortUrl(ctx, keeper)
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, totalCost)
		if err != nil {
			rullBackNumber()
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
		keeper.StoreLAddress(ctx, newSUrl, msg.Buyer, msg.Bid, DefaultRentDuration)
		GlobalBloomFilter.Set(newSUrl)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeBuySUrl,
			sdk.NewAttribute(AttributeKeySUrl, msg.SUrl),
			sdk.NewAttribute(AttributeKeyBuyer, msg.Buyer.String()),
			sdk.NewAttribute(AttributeKeyPrice, msg.Bid.String()),
			sdk.NewAttribute(AttributeKeyExpiration, time.Now().Add(DefaultRentDuration).Format(time.RFC3339)),
		),
	)

	return sdk.Result{}
}

func txUrl(ctx sdk.Context, keeper Keeper, msg MsgBuySUrl, fee sdk.Coins) sdk.Result {
	if keeper.HasOwner(ctx, msg.SUrl) && !keeper.GetSell(ctx, msg.SUrl) {
		return sdk.ErrInternal("Address does not sell").Result()
	}
	if keeper.IsExpired(ctx, msg.SUrl) {
		keeper.DeleteLAddress(ctx, msg.SUrl)
		return sdk.ErrUnauthorized("SUrl has expired").Result()
	}
	if keeper.GetPrice(ctx, msg.SUrl).IsAllGT(msg.Bid) {
		return sdk.ErrInsufficientCoins("Bid not high enough").Result()
	}

	totalCost := msg.Bid.Add(fee)
	_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.SUrl), msg.Bid)
	if err != nil {
		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	}

	keeper.SetOwner(ctx, msg.SUrl, msg.Buyer)
	keeper.SetPrice(ctx, msg.SUrl, msg.Bid)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeBuySUrl,
			sdk.NewAttribute(AttributeKeySUrl, msg.SUrl),
			sdk.NewAttribute(AttributeKeyBuyer, msg.Buyer.String()),
			sdk.NewAttribute(AttributeKeyOwner, keeper.GetOwner(ctx, msg.SUrl).String()),
			sdk.NewAttribute(AttributeKeyPrice, msg.Bid.String()),
		),
	)

	return sdk.Result{}
}

func handleMsgSetSell(ctx sdk.Context, keeper Keeper, msg MsgSetSell) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl) {
		return sdk.ErrUnauthorized("Incorrect Surl").Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	if keeper.IsExpired(ctx, msg.SUrl) {
		return sdk.ErrUnauthorized("SUrl has expired").Result()
	}
	keeper.SetSell(ctx, msg.SUrl, msg.IsSell)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeSetSell,
			sdk.NewAttribute(AttributeKeySUrl, msg.SUrl),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(AttributeKeyIsSell, fmt.Sprintf("%t", msg.IsSell)),
		),
	)

	return sdk.Result{}
}

func handleMsgSetPrice(ctx sdk.Context, keeper Keeper, msg MsgSetPrice) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl) {
		return sdk.ErrUnauthorized("Incorrect Surl").Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	if keeper.IsExpired(ctx, msg.SUrl) {
		return sdk.ErrUnauthorized("SUrl has expired").Result()
	}
	keeper.SetPrice(ctx, msg.SUrl, msg.Bid)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeSetPrice,
			sdk.NewAttribute(AttributeKeySUrl, msg.SUrl),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(AttributeKeyPrice, msg.Bid.String()),
		),
	)

	return sdk.Result{}
}

func handleMsgRenew(ctx sdk.Context, keeper Keeper, msg MsgRenew) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl) {
		return sdk.ErrUnauthorized("Incorrect Surl").Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}

	lAddress := keeper.GetLAddress(ctx, msg.SUrl)
	existingExp := lAddress.ExpirationTime
	if existingExp.Before(time.Now()) {
		existingExp = time.Now()
	}
	newExp := existingExp.Add(msg.Duration)

	fee := calculateRenewFee(msg.Duration)
	totalCost := fee

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, totalCost)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough coins").Result()
	}

	keeper.Renew(ctx, msg.SUrl, msg.Duration)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeRenew,
			sdk.NewAttribute(AttributeKeySUrl, msg.SUrl),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(AttributeKeyExpiration, newExp.Format(time.RFC3339)),
		),
	)

	return sdk.Result{}
}

func calculateFee(bid sdk.Coins) sdk.Coins {
	var fee sdk.Coins
	for _, coin := range bid {
		fee = fee.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, coin.Amount.Mul(sdk.NewInt(5)).Div(sdk.NewInt(100))))
	}
	return fee
}

func calculateRenewFee(duration time.Duration) sdk.Coins {
	days := int(duration.Hour() / 24)
	feeAmount := sdk.NewInt(int64(days))
	return sdk.NewCoins(sdk.NewCoin("sastoken", feeAmount))
}
