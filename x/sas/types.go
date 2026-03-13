package sas

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	time "time"
)

const (
	MaxLUrlLength       = 2048
	DefaultRentDuration = 365 * 24 * time.Hour
)

var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("sastoken", 1)}

type LAddress struct {
	LUrl           string         `json:"lUrl"`
	Owner          sdk.AccAddress `json:"owner"`
	Price          sdk.Coins      `json:"price"`
	IsSell         bool           `json:"isSell"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	ExpirationTime time.Time      `json:"expirationTime"`
}

func NewLAddress() LAddress {
	return LAddress{
		Price:     MinNamePrice,
		CreatedAt: time.Now(),
	}
}
