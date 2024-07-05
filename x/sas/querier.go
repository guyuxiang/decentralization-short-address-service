package sas

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the sas Querier
const (
	QueryLUrl    = "LUrl"
	QueryLAdress = "LAdress"
	QuerySUrl    = "SUrls"
)

// Query Result Payload for a resolve query
type QueryResLUrl struct {
	Lurl string `json:"lUrl"`
}

// implement fmt.Stringer
func (r QueryResLUrl) String() string {
	return r.Lurl
}

// Query Result Payload for a names query
type QueryResSUrls []string

// implement fmt.Stringer
func (n QueryResSUrls) String() string {
	return strings.Join(n[:], "\n")
}

// implement fmt.Stringer
func (w LAdress) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
LUrl: %s
Price: %s
IsSell: %s`, w.Owner, w.LUrl, w.Price, w.IsSell))
}

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryLUrl:
			return queryLUrl(ctx, path[1:], req, keeper)
		case QueryLAdress:
			return queryLAdress(ctx, path[1:], req, keeper)
		case QuerySUrl:
			return querySUrls(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown sas query endpoint")
		}
	}
}

// nolint: unparam
func queryLUrl(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	sUrl := path[0]
	if !QuickCheckSUrlExist(sUrl) {
		return nil, sdk.ErrUnknownRequest("sUrl not exist")
	}
	// 缓存LRU
	lUrl, fit := LruCache.Get(sUrl)
	if !fit {
		lu := keeper.ResolveLUrl(ctx, sUrl)
		if lu == "" {
			return []byte{}, sdk.ErrUnknownRequest("could not resolve name")
		}
		LruCache.Set(sUrl, lu)
		lUrl = lu
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, QueryResLUrl{lUrl.(string)})
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

// nolint: unparam
func queryLAdress(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	sUrl := path[0]
	if !QuickCheckSUrlExist(sUrl) {
		return nil, sdk.ErrUnknownRequest("sUrl not exist")
	}
	lAdress := keeper.GetLAdress(ctx, sUrl)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, lAdress)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func querySUrls(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var sUrlsList QueryResSUrls

	iterator := keeper.GetSUrlsIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		sUrl := string(iterator.Key())
		sUrlsList = append(sUrlsList, sUrl)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, sUrlsList)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}
