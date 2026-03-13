package sas

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey

	cdc *codec.Codec
}

func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) GetLAddress(ctx sdk.Context, sUrl string) LAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(sUrl))
	var lAddress LAddress
	k.cdc.MustUnmarshalBinaryBare(bz, &lAddress)
	return lAddress
}

func (k Keeper) StoreLAddress(ctx sdk.Context, sUrl string, owner sdk.AccAddress, price sdk.Coins, duration time.Duration) {
	lAddress := NewLAddress()
	lAddress.Price = price
	lAddress.Owner = owner
	lAddress.ExpirationTime = time.Now().Add(duration)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(sUrl), k.cdc.MustMarshalBinaryBare(lAddress))
}

func (k Keeper) isSUrlExist(ctx sdk.Context, sUrl string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(sUrl))
}

func (k Keeper) SetLAddress(ctx sdk.Context, sUrl string, lAddress LAddress) {
	if lAddress.Owner.Empty() {
		return
	}
	lAddress.UpdatedAt = time.Now()
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(sUrl), k.cdc.MustMarshalBinaryBare(lAddress))
}

func (k Keeper) ResolveLUrl(ctx sdk.Context, sUrl string) string {
	return k.GetLAddress(ctx, sUrl).LUrl
}

func (k Keeper) SetLUrl(ctx sdk.Context, sUrl string, lUrl string) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.LUrl = lUrl
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) HasOwner(ctx sdk.Context, sUrl string) bool {
	return !k.GetLAddress(ctx, sUrl).Owner.Empty()
}

func (k Keeper) GetOwner(ctx sdk.Context, sUrl string) sdk.AccAddress {
	return k.GetLAddress(ctx, sUrl).Owner
}

func (k Keeper) SetOwner(ctx sdk.Context, sUrl string, owner sdk.AccAddress) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.Owner = owner
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) GetPrice(ctx sdk.Context, sUrl string) sdk.Coins {
	return k.GetLAddress(ctx, sUrl).Price
}

func (k Keeper) SetPrice(ctx sdk.Context, sUrl string, price sdk.Coins) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.Price = price
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) GetSUrlsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

func (k Keeper) SetSell(ctx sdk.Context, sUrl string, isSell bool) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.IsSell = isSell
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) SetNoSell(ctx sdk.Context, sUrl string) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.IsSell = false
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) GetSell(ctx sdk.Context, sUrl string) bool {
	return k.GetLAddress(ctx, sUrl).IsSell
}

func (k Keeper) IsExpired(ctx sdk.Context, sUrl string) bool {
	lAddress := k.GetLAddress(ctx, sUrl)
	return time.Now().After(lAddress.ExpirationTime)
}

func (k Keeper) GetExpirationTime(ctx sdk.Context, sUrl string) time.Time {
	return k.GetLAddress(ctx, sUrl).ExpirationTime
}

func (k Keeper) Renew(ctx sdk.Context, sUrl string, duration time.Duration) sdk.Error {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.ExpirationTime = time.Now().Add(duration)
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
	return nil
}

func (k Keeper) DeleteLAddress(ctx sdk.Context, sUrl string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(sUrl))
}
