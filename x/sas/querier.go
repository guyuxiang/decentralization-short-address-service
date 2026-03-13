package sas

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryLUrl     = "LUrl"
	QueryLAddress = "LAddress"
	QuerySUrl     = "SUrls"
	QueryOwner    = "owner"
	QueryStats    = "stats"
	DefaultPage   = 1
	DefaultLimit  = 100
)

type QueryResLUrl struct {
	LUrl string `json:"lUrl"`
}

func (r QueryResLUrl) String() string {
	return r.LUrl
}

type QueryResSUrls []string

func (n QueryResSUrls) String() string {
	return strings.Join(n[:], "\n")
}

type QueryResPage struct {
	SUrls []string `json:"sUrls"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Total int      `json:"total"`
}

func (w LAddress) String() string {
	return fmt.Sprintf(`Owner: %s
LUrl: %s
Price: %s
IsSell: %v
CreatedAt: %s
UpdatedAt: %s
ExpirationTime: %s`, w.Owner, w.LUrl, w.Price, w.IsSell, w.CreatedAt, w.UpdatedAt, w.ExpirationTime)
}

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryLUrl:
			return queryLUrl(ctx, path[1:], req, keeper)
		case QueryLAddress:
			return queryLAddress(ctx, path[1:], req, keeper)
		case QuerySUrl:
			return querySUrls(ctx, req, keeper)
		case QueryOwner:
			return queryOwnerSUrls(ctx, path[1:], req, keeper)
		case QueryStats:
			return queryStats(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown sas query endpoint")
		}
	}
}

func queryLUrl(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	sUrl := path[0]
	if !QuickCheckSUrlExist(sUrl) {
		return nil, sdk.ErrUnknownRequest("sUrl not exist")
	}
	lUrl, fit := LruCache.Get(sUrl)
	if !fit {
		lu := keeper.ResolveLUrl(ctx, sUrl)
		if lu == "" {
			return []byte{}, sdk.ErrUnknownRequest("could not resolve name")
		}
		LruCache.Set(sUrl, lu)
		lUrl = lu
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, QueryResLUrl{LUrl: lUrl.(string)})
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryLAddress(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	sUrl := path[0]
	if !QuickCheckSUrlExist(sUrl) {
		return nil, sdk.ErrUnknownRequest("sUrl not exist")
	}
	lAddress := keeper.GetLAddress(ctx, sUrl)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, lAddress)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func querySUrls(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	page := DefaultPage
	limit := DefaultLimit

	if req.Data != nil && len(req.Data) > 0 {
		var params struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
		}
		if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err == nil {
			if params.Page > 0 {
				page = params.Page
			}
			if params.Limit > 0 {
				limit = params.Limit
			}
		}
	}

	var allSUrls []string
	iterator := keeper.GetSUrlsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		sUrl := string(iterator.Key())
		allSUrls = append(allSUrls, sUrl)
	}

	total := len(allSUrls)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		allSUrls = []string{}
	} else {
		if end > total {
			end = total
		}
		allSUrls = allSUrls[start:end]
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, QueryResPage{
		SUrls: allSUrls,
		Page:  page,
		Limit: limit,
		Total: total,
	})
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryOwnerSUrls(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	if len(path) < 1 {
		return nil, sdk.ErrUnknownRequest("owner address required")
	}

	ownerStr := path[0]
	owner, err := sdk.AccAddressFromBech32(ownerStr)
	if err != nil {
		return nil, sdk.ErrInvalidAddress(ownerStr)
	}

	page := DefaultPage
	limit := DefaultLimit

	if len(path) > 1 {
		fmt.Sscanf(path[1], "%d", &page)
	}
	if len(path) > 2 {
		fmt.Sscanf(path[2], "%d", &limit)
	}

	var ownerSUrls []string
	iterator := keeper.GetSUrlsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		sUrl := string(iterator.Key())
		lAddress := keeper.GetLAddress(ctx, sUrl)
		if lAddress.Owner.Equals(owner) {
			ownerSUrls = append(ownerSUrls, sUrl)
		}
	}

	total := len(ownerSUrls)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		ownerSUrls = []string{}
	} else {
		if end > total {
			end = total
		}
		ownerSUrls = ownerSUrls[start:end]
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, QueryResPage{
		SUrls: ownerSUrls,
		Page:  page,
		Limit: limit,
		Total: total,
	})
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryStats(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	stats := keeper.GetStats(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, stats)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}
