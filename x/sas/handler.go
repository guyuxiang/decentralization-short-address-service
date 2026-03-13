package sas

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		case MsgBuySUrlEscrow:
			return handleMsgBuySUrlEscrow(ctx, keeper, msg)
		case MsgConfirmEscrow:
			return handleMsgConfirmEscrow(ctx, keeper, msg)
		case MsgCancelEscrow:
			return handleMsgCancelEscrow(ctx, keeper, msg)
		case MsgBatchSetLUrl:
			return handleMsgBatchSetLUrl(ctx, keeper, msg)
		case MsgAddBlackList:
			return handleMsgAddBlackList(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized sas Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSetLUrl(ctx sdk.Context, keeper Keeper, msg MsgSetLUrl) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl, len(msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Surl").Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	if keeper.IsExpired(ctx, msg.SUrl) {
		return sdk.ErrUnauthorized("SUrl has expired").Result()
	}
	if keeper.IsBlackListed(msg.LUrl) {
		return sdk.ErrUnauthorized("URL is blacklisted").Result()
	}

	keeper.SetLUrl(ctx, msg.SUrl, msg.LUrl)
	return sdk.Result{}
}

func handleMsgBuySUrl(ctx sdk.Context, keeper Keeper, msg MsgBuySUrl) sdk.Result {
	fee := calculateFee(msg.Bid)

	urlLength := msg.Length
	if urlLength == 0 {
		urlLength = 6
	}

	if len(msg.SUrl) != 0 {
		if CheckSUrlExist(ctx, keeper, msg.SUrl, len(msg.SUrl)) {
			if keeper.IsExpired(ctx, msg.SUrl) {
				keeper.DeleteLAddress(ctx, msg.SUrl)
			} else {
				return txUrl(ctx, keeper, msg, fee)
			}
		}
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins for bid").Result()
		}
		_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, fee)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins for fee").Result()
		}
		keeper.StoreLAddress(ctx, msg.SUrl, msg.Buyer, msg.Bid, DefaultRentDuration)
		keeper.AddToBloomFilter(msg.SUrl)
	} else {
		newSUrl := ApplyShortUrl(ctx, keeper, urlLength)
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid)
		if err != nil {
			rullBackNumber()
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins for bid").Result()
		}
		_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, fee)
		if err != nil {
			rullBackNumber()
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins for fee").Result()
		}
		keeper.StoreLAddress(ctx, newSUrl, msg.Buyer, msg.Bid, DefaultRentDuration)
		keeper.AddToBloomFilter(newSUrl)
	}

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

	_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.SUrl), msg.Bid)
	if err != nil {
		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	}

	keeper.SetOwner(ctx, msg.SUrl, msg.Buyer)
	keeper.SetPrice(ctx, msg.SUrl, msg.Bid)
	return sdk.Result{}
}

func handleMsgSetSell(ctx sdk.Context, keeper Keeper, msg MsgSetSell) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl, len(msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Surl").Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	if keeper.IsExpired(ctx, msg.SUrl) {
		return sdk.ErrUnauthorized("SUrl has expired").Result()
	}
	keeper.SetSell(ctx, msg.SUrl, msg.IsSell)
	return sdk.Result{}
}

func handleMsgSetPrice(ctx sdk.Context, keeper Keeper, msg MsgSetPrice) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl, len(msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Surl").Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	if keeper.IsExpired(ctx, msg.SUrl) {
		return sdk.ErrUnauthorized("SUrl has expired").Result()
	}
	keeper.SetPrice(ctx, msg.SUrl, msg.Bid)
	return sdk.Result{}
}

func handleMsgRenew(ctx sdk.Context, keeper Keeper, msg MsgRenew) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl, len(msg.SUrl)) {
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

	fee := calculateRenewFee(msg.Duration)
	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, fee)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough coins").Result()
	}

	keeper.Renew(ctx, msg.SUrl, msg.Duration)
	return sdk.Result{}
}

func handleMsgBuySUrlEscrow(ctx sdk.Context, keeper Keeper, msg MsgBuySUrlEscrow) sdk.Result {
	if !CheckSUrlExist(ctx, keeper, msg.SUrl, len(msg.SUrl)) {
		return sdk.ErrUnauthorized("SUrl does not exist").Result()
	}
	if keeper.IsExpired(ctx, msg.SUrl) {
		return sdk.ErrUnauthorized("SUrl has expired").Result()
	}
	if !keeper.GetSell(ctx, msg.SUrl) {
		return sdk.ErrUnauthorized("Address is not for sale").Result()
	}

	fee := calculateFee(msg.Amount)
	totalCost := msg.Amount.Add(fee)

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, totalCost)
	if err != nil {
		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	}

	escrow := NewEscrow(msg.SUrl, keeper.GetOwner(ctx, msg.SUrl), msg.Buyer, msg.Amount, fee)
	store := ctx.KVStore(keeper.storeKey)
	store.Set([]byte("escrow_"+msg.SUrl), keeper.cdc.MustMarshalJSON(escrow))

	return sdk.Result{}
}

func handleMsgConfirmEscrow(ctx sdk.Context, keeper Keeper, msg MsgConfirmEscrow) sdk.Result {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get([]byte("escrow_" + msg.SUrl))
	if bz == nil {
		return sdk.ErrInternal("Escrow not found").Result()
	}

	var escrow Escrow
	keeper.cdc.MustUnmarshalJSON(bz, &escrow)

	if escrow.Status != EscrowPending {
		return sdk.ErrInternal("Escrow is not pending").Result()
	}
	if !msg.Confirmor.Equals(escrow.Seller) && !msg.Confirmor.Equals(escrow.Buyer) {
		return sdk.ErrUnauthorized("Not authorized to confirm").Result()
	}

	_, err := keeper.coinKeeper.SendCoins(ctx, escrow.Buyer, escrow.Seller, escrow.Amount)
	if err != nil {
		keeper.coinKeeper.AddCoins(ctx, escrow.Buyer, escrow.Amount.Add(escrow.Fee))
		return sdk.ErrInsufficientCoins("Transaction failed").Result()
	}

	now := time.Now()
	escrow.Status = EscrowCompleted
	escrow.CompletedAt = &now
	store.Set([]byte("escrow_"+msg.SUrl), keeper.cdc.MustMarshalJSON(escrow))

	keeper.SetOwner(ctx, msg.SUrl, escrow.Buyer)
	return sdk.Result{}
}

func handleMsgCancelEscrow(ctx sdk.Context, keeper Keeper, msg MsgCancelEscrow) sdk.Result {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get([]byte("escrow_" + msg.SUrl))
	if bz == nil {
		return sdk.ErrInternal("Escrow not found").Result()
	}

	var escrow Escrow
	keeper.cdc.MustUnmarshalJSON(bz, &escrow)

	if escrow.Status != EscrowPending {
		return sdk.ErrInternal("Escrow is not pending").Result()
	}
	if !msg.Canceler.Equals(escrow.Seller) && !msg.Canceler.Equals(escrow.Buyer) {
		return sdk.ErrUnauthorized("Not authorized to cancel").Result()
	}

	keeper.coinKeeper.AddCoins(ctx, escrow.Buyer, escrow.Amount.Add(escrow.Fee))

	escrow.Status = EscrowCancelled
	store.Set([]byte("escrow_"+msg.SUrl), keeper.cdc.MustMarshalJSON(escrow))

	return sdk.Result{}
}

func handleMsgBatchSetLUrl(ctx sdk.Context, keeper Keeper, msg MsgBatchSetLUrl) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.SUrl)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	if keeper.IsExpired(ctx, msg.SUrl) {
		return sdk.ErrUnauthorized("SUrl has expired").Result()
	}

	successCount := 0
	for _, lUrl := range msg.LUrls {
		if keeper.IsBlackListed(lUrl) {
			continue
		}
		keeper.SetLUrl(ctx, msg.SUrl, lUrl)
		successCount++
	}

	if successCount == 0 && len(msg.LUrls) > 0 {
		return sdk.ErrUnauthorized("All URLs are blacklisted").Result()
	}

	return sdk.Result{}
}

func handleMsgAddBlackList(ctx sdk.Context, keeper Keeper, msg MsgAddBlackList) sdk.Result {
	if msg.IsDomain {
		keeper.AddToBlackListDomain(ctx, msg.URL)
	} else {
		keeper.AddToBlackListURL(ctx, msg.URL)
	}
	return sdk.Result{}
}

func calculateFee(bid sdk.Coins) sdk.Coins {
	totalFee := sdk.NewInt(0)
	for _, coin := range bid {
		fee := coin.Amount.Mul(sdk.NewInt(5)).Quo(sdk.NewInt(100))
		totalFee = totalFee.Add(fee)
	}
	if totalFee.IsZero() {
		totalFee = sdk.NewInt(1)
	}
	return sdk.Coins{{Denom: "stake", Amount: totalFee}}
}

func calculateRenewFee(duration time.Duration) sdk.Coins {
	days := int(duration.Hours() / 24)
	if days < 1 {
		days = 1
	}
	return sdk.Coins{{Denom: "stake", Amount: sdk.NewInt(int64(days))}}
}
