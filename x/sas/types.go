package sas

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("sastoken", 1)}

// Whois is a struct that contains all the metadata of a name
type LAdress struct {
	LUrl   string         `json:"lUrl"`
	Owner  sdk.AccAddress `json:"owner"`
	Price  sdk.Coins      `json:"price"`
	IsSell bool           `json:"isSell"`
}

// Returns a new Whois with the minprice as the price
func NewLAdress() LAdress {
	return LAdress{
		Price: MinNamePrice,
	}
}
