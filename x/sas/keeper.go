package sas

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the sas Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) GetLAdress(ctx sdk.Context, sUrl string) LAdress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(sUrl))
	var lAdress LAdress
	k.cdc.MustUnmarshalBinaryBare(bz, &lAdress)
	return lAdress
}

func (k Keeper) StoreLAdress(ctx sdk.Context, sUrl string, owner sdk.AccAddress, price sdk.Coins) {
	lAdress := NewLAdress()
	lAdress.Price = price
	lAdress.Owner = owner
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(sUrl), k.cdc.MustMarshalBinaryBare(lAdress))
}

func (k Keeper) isSUrlExist(ctx sdk.Context, sUrl string) bool {
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(sUrl)) {
		return true
	} else {
		return false
	}
}

func (k Keeper) SetLAdress(ctx sdk.Context, sUrl string, lAdress LAdress) {
	if lAdress.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(sUrl), k.cdc.MustMarshalBinaryBare(lAdress))
}

func (k Keeper) ResolveLUrl(ctx sdk.Context, sUrl string) string {
	return k.GetLAdress(ctx, sUrl).LUrl
}

func (k Keeper) SetLUrl(ctx sdk.Context, sUrl string, lUrl string) {
	lAdress := k.GetLAdress(ctx, sUrl)
	lAdress.LUrl = lUrl
	k.SetLAdress(ctx, sUrl, lAdress)
}

func (k Keeper) HasOwner(ctx sdk.Context, sUrl string) bool {
	return !k.GetLAdress(ctx, sUrl).Owner.Empty()
}

func (k Keeper) GetOwner(ctx sdk.Context, sUrl string) sdk.AccAddress {
	return k.GetLAdress(ctx, sUrl).Owner
}

func (k Keeper) SetOwner(ctx sdk.Context, sUrl string, owner sdk.AccAddress) {
	lAdress := k.GetLAdress(ctx, sUrl)
	lAdress.Owner = owner
	k.SetLAdress(ctx, sUrl, lAdress)
}

func (k Keeper) GetPrice(ctx sdk.Context, sUrl string) sdk.Coins {
	return k.GetLAdress(ctx, sUrl).Price
}

func (k Keeper) SetPrice(ctx sdk.Context, sUrl string, price sdk.Coins) {
	lAdress := k.GetLAdress(ctx, sUrl)
	lAdress.Price = price
	k.SetLAdress(ctx, sUrl, lAdress)
}

func (k Keeper) GetSUrlsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

func (k Keeper) SetSell(ctx sdk.Context, sUrl string, isSell bool) {
	lAdress := k.GetLAdress(ctx, sUrl)
	lAdress.IsSell = isSell
	k.SetLAdress(ctx, sUrl, lAdress)
}

func (k Keeper) SetNoSell(ctx sdk.Context, sUrl string) {
	lAdress := k.GetLAdress(ctx, sUrl)
	lAdress.IsSell = false
	k.SetLAdress(ctx, sUrl, lAdress)
}

func (k Keeper) GetSell(ctx sdk.Context, sUrl string) bool {
	return k.GetLAdress(ctx, sUrl).IsSell
}
