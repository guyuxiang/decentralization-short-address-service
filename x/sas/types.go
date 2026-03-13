package sas

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MaxLUrlLength       = 2048
	DefaultRentDuration = 365 * 24 * time.Hour
	GracePeriod         = 7 * 24 * time.Hour
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
	ClickCount     uint64         `json:"clickCount"`
}

func NewLAddress() LAddress {
	return LAddress{
		Price:     MinNamePrice,
		CreatedAt: time.Now(),
	}
}

type BlackList struct {
	URLs      map[string]bool `json:"urls"`
	Domains   map[string]bool `json:"domains"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

func NewBlackList() BlackList {
	return BlackList{
		URLs:      make(map[string]bool),
		Domains:   make(map[string]bool),
		UpdatedAt: time.Now(),
	}
}

type Stats struct {
	TotalClicks  uint64            `json:"totalClicks"`
	DailyClicks  map[string]uint64 `json:"dailyClicks"`
	TopShortURLs []StatEntry       `json:"topShortUrls"`
	TopOwners    []StatEntry       `json:"topOwners"`
}

type StatEntry struct {
	Key   string `json:"key"`
	Value uint64 `json:"value"`
}

func NewStats() Stats {
	return Stats{
		TotalClicks:  0,
		DailyClicks:  make(map[string]uint64),
		TopShortURLs: make([]StatEntry, 0),
		TopOwners:    make([]StatEntry, 0),
	}
}

type Escrow struct {
	SUrl        string         `json:"sUrl"`
	Seller      sdk.AccAddress `json:"seller"`
	Buyer       sdk.AccAddress `json:"buyer"`
	Amount      sdk.Coins      `json:"amount"`
	Fee         sdk.Coins      `json:"fee"`
	Status      EscrowStatus   `json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	CompletedAt *time.Time     `json:"completedAt,omitempty"`
}

type EscrowStatus int

const (
	EscrowPending EscrowStatus = iota
	EscrowCompleted
	EscrowCancelled
)

func NewEscrow(sUrl string, seller, buyer sdk.AccAddress, amount, fee sdk.Coins) Escrow {
	return Escrow{
		SUrl:      sUrl,
		Seller:    seller,
		Buyer:     buyer,
		Amount:    amount,
		Fee:       fee,
		Status:    EscrowPending,
		CreatedAt: time.Now(),
	}
}

type ExpiredSUrl struct {
	SUrl        string    `json:"sUrl"`
	ExpiredAt   time.Time `json:"expiredAt"`
	OriginalExp time.Time `json:"originalExp"`
}
