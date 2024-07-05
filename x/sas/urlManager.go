package sas

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sync/atomic"
)

var gC = &globeCounter{}

type globeCounter struct {
	number *uint32
}

func NewGlobeCounter(startNumber uint32) {
	atomic.AddUint32(gC.number, startNumber)
}

func GetNumber() uint32 {
	return *gC.number
}

func applyNumber() uint32 {
	atomic.AddUint32(gC.number, 1)
	return *gC.number
}

func rullBackNumber() {
	if *gC.number >= 1 {
		atomic.CompareAndSwapUint32(gC.number, *gC.number, *gC.number-1)
	}
}

func ApplyShortUrl(ctx sdk.Context, keeper Keeper) string {
	number := applyNumber()
	newSUrl := Encode6(int(number))
	if CheckSUrlExist(ctx, keeper, newSUrl) {
		newSUrl = ApplyShortUrl(ctx, keeper)
	}
	return newSUrl
}

func CheckSUrlExist(ctx sdk.Context, keeper Keeper, sUrl string) bool {
	if !QuickCheckSUrlExist(sUrl) {
		return false
	}
	if Decode(sUrl) <= int(GetNumber()) || keeper.isSUrlExist(ctx, sUrl) {
		return true
	}
	return false
}

func QuickCheckSUrlExist(sUrl string) bool {
	if GlobalBloomFilter.Check(sUrl) == false {
		return false
	}
	return true
}
