package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"sas/x/sas"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

const (
	restName = "sUrl"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/adress/sUrls", storeName), sUrlsHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/adress", storeName), buySUrlHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/adress/lUrl", storeName), setLUrlHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/adress/price", storeName), setPriceHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/adress/sell", storeName), setSellHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/adress/{%s}/lUrl", storeName, restName), lUrlHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/adress/{%s}/lAddress", storeName, restName), lAddressHandler(cdc, cliCtx, storeName)).Methods("GET")

	// Redirect route -访问短地址自动跳转长地址
	r.HandleFunc("/s/{sUrl}", redirectHandler(cliCtx, storeName)).Methods("GET")

	// Stats route
	r.HandleFunc(fmt.Sprintf("/%s/stats", storeName), statsHandler(cdc, cliCtx, storeName)).Methods("GET")
}

type buySUrlReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	SUrl    string       `json:"sUrl"`
	Amount  string       `json:"amount"`
	Buyer   string       `json:"buyer"`
}

func buySUrlHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req buySUrlReq

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Buyer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		coins, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := sas.NewMsgBuySUrl(req.SUrl, coins, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type setLUrlReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	SUrl    string       `json:"sUrl"`
	LUrl    string       `json:"lUrl"`
	Owner   string       `json:"owner"`
}

func setLUrlHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setLUrlReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := sas.NewMsgSetLUrl(req.SUrl, req.LUrl, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type setPriceReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	SUrl    string       `json:"sUrl"`
	Amount  string       `json:"amount"`
	Owner   string       `json:"owner"`
}

func setPriceHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setPriceReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		coins, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := sas.NewMsgSetPrice(req.SUrl, coins, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func lUrlHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/LUrl/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

type setSellReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	SUrl    string       `json:"sUrl"`
	IsSell  bool         `json:"isSell"`
	Owner   string       `json:"owner"`
}

func setSellHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setSellReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := sas.NewMsgSetSell(req.SUrl, req.IsSell, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func lAddressHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/LAddress/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func sUrlsHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/SUrls", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func redirectHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sUrl := vars["sUrl"]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/LUrl/%s", storeName, sUrl), nil)
		if err != nil || len(res) == 0 {
			http.NotFound(w, r)
			return
		}

		var lUrlResp struct {
			LUrl string `json:"lUrl"`
		}
		if err := cdc.UnmarshalJSON(res, &lUrlResp); err != nil || lUrlResp.LUrl == "" {
			http.NotFound(w, r)
			return
		}

		lUrl := strings.ToLower(lUrlResp.LUrl)
		if !strings.HasPrefix(lUrl, "http://") && !strings.HasPrefix(lUrl, "https://") {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, lUrlResp.LUrl, http.StatusMovedPermanently)
	}
}

func statsHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/stats", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
