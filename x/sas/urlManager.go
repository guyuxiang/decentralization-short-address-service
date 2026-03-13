package sas

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sync/atomic"
)

var gC = &globeCounter{
	number: new(uint32),
}

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
	atomic.AddUint32(gC.number, ^uint32(0))
}

func ApplyShortUrl(ctx sdk.Context, keeper Keeper, length int) string {
	if length < 1 {
		length = 6
	}
	if length > 6 {
		length = 6
	}
	number := applyNumber()
	newSUrl := EncodeFixedLength(int(number), length)
	if CheckSUrlExist(ctx, keeper, newSUrl, length) {
		newSUrl = ApplyShortUrl(ctx, keeper, length)
	}
	return newSUrl
}

func CheckSUrlExist(ctx sdk.Context, keeper Keeper, sUrl string, length int) bool {
	if !QuickCheckSUrlExist(sUrl) {
		return false
	}
	maxNum := calculateMaxNumberForLength(length)
	if Decode(sUrl) <= int(GetNumber()) && Decode(sUrl) <= maxNum || keeper.isSUrlExist(ctx, sUrl) {
		return true
	}
	return false
}

func calculateMaxNumberForLength(length int) int {
	result := 1
	for i := 0; i < length; i++ {
		result = result * 62
	}
	return result - 1
}

func QuickCheckSUrlExist(sUrl string) bool {
	if GlobalBloomFilter.Check(sUrl) == false {
		return false
	}
	return true
}

func (k Keeper) AddToBloomFilter(sUrl string) {
	if GlobalBloomFilter != nil {
		GlobalBloomFilter.Set(sUrl)
	}
}

func (k Keeper) RemoveFromBloomFilter(sUrl string) {
}
