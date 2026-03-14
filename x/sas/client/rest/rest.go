package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"sas/x/sas"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
	
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	restName = "sUrl"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/adress/sUrls", storeName), sUrlsHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/adress/sUrls/detail", storeName), sUrlsDetailHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/adress/owner/{owner}", storeName), ownerUrlsHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/adress", storeName), buySUrlHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/adress/lUrl", storeName), setLUrlHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/adress/price", storeName), setPriceHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/adress/sell", storeName), setSellHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/adress/{%s}/lUrl", storeName, restName), lUrlHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/adress/{%s}/lAddress", storeName, restName), lAddressHandler(cdc, cliCtx, storeName)).Methods("GET")

	// Faucet route -领取测试代币
	r.HandleFunc(fmt.Sprintf("/%s/faucet", storeName), faucetHandler(cdc, cliCtx)).Methods("POST")

	// Stats route
	r.HandleFunc(fmt.Sprintf("/%s/stats", storeName), statsHandler(cdc, cliCtx, storeName)).Methods("GET")

	// Redirect route -访问短地址自动跳转长地址
	r.HandleFunc("/{sUrl}", redirectHandler(cliCtx, storeName, cdc)).Methods("GET")
}

type buySUrlReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	SUrl    string       `json:"sUrl"`
	Amount  string       `json:"amount"`
	Buyer   string       `json:"buyer"`
	Length  int          `json:"length"`
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
		msg := sas.NewMsgBuySUrl(req.SUrl, coins, addr, req.Length)
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

type SUrlDetail struct {
	SUrl          string `json:"sUrl"`
	LUrl          string `json:"lUrl"`
	Owner         string `json:"owner"`
	Price         string `json:"price"`
	IsSell        bool   `json:"isSell"`
	ExpirationTime string `json:"expirationTime"`
	Clicks        uint64 `json:"clicks"`
}

type SUrlDetailResponse struct {
	Result []SUrlDetail `json:"result"`
}

func sUrlsDetailHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/SUrls", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var pageRes struct {
			SUrls []string `json:"sUrls"`
		}
		if err := cdc.UnmarshalJSON(res, &pageRes); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var details []SUrlDetail
		for _, sUrl := range pageRes.SUrls {
			lAddrRes, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/LAddress/%s", storeName, sUrl), nil)
			if err != nil || len(lAddrRes) == 0 {
				continue
			}

			var lAddr struct {
				LUrl           string `json:"lUrl"`
				Owner          string `json:"owner"`
				Price          string `json:"price"`
				IsSell         bool   `json:"isSell"`
				ExpirationTime string `json:"expirationTime"`
				ClickCount     uint64 `json:"clickCount"`
			}
			if err := cdc.UnmarshalJSON(lAddrRes, &lAddr); err != nil {
				continue
			}

			details = append(details, SUrlDetail{
				SUrl:            sUrl,
				LUrl:            lAddr.LUrl,
				Owner:           lAddr.Owner,
				Price:           lAddr.Price,
				IsSell:          lAddr.IsSell,
				ExpirationTime:  lAddr.ExpirationTime,
				Clicks:          lAddr.ClickCount,
			})
		}

		resp := SUrlDetailResponse{Result: details}
		rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
	}
}

func ownerUrlsHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		owner := vars["owner"]

		if _, err := sdk.AccAddressFromBech32(owner); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid owner address")
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/SUrls", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var pageRes struct {
			SUrls []string `json:"sUrls"`
		}
		if err := cdc.UnmarshalJSON(res, &pageRes); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var details []SUrlDetail
		for _, sUrl := range pageRes.SUrls {
			lAddrRes, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/LAddress/%s", storeName, sUrl), nil)
			if err != nil || len(lAddrRes) == 0 {
				continue
			}

			var lAddr struct {
				LUrl           string `json:"lUrl"`
				Owner          string `json:"owner"`
				Price          string `json:"price"`
				IsSell         bool   `json:"isSell"`
				ExpirationTime string `json:"expirationTime"`
				ClickCount     uint64 `json:"clickCount"`
			}
			if err := cdc.UnmarshalJSON(lAddrRes, &lAddr); err != nil {
				continue
			}

			if lAddr.Owner == owner {
				details = append(details, SUrlDetail{
					SUrl:            sUrl,
					LUrl:            lAddr.LUrl,
					Owner:           lAddr.Owner,
					Price:           lAddr.Price,
					IsSell:          lAddr.IsSell,
					ExpirationTime:  lAddr.ExpirationTime,
					Clicks:          lAddr.ClickCount,
				})
			}
		}

		resp := SUrlDetailResponse{Result: details}
		rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
	}
}

func redirectHandler(cliCtx context.CLIContext, storeName string, cdc *codec.Codec) http.HandlerFunc {
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

type faucetReq struct {
	Address string `json:"address"`
}

type faucetResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	TxHash  string `json:"tx_hash,omitempty"`
}

func faucetHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	faucetAmount := "1000os"
	
	return func(w http.ResponseWriter, r *http.Request) {
		var req faucetReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if req.Address == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "address is required")
			return
		}

		toAddr, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid address format")
			return
		}

		faucetAddr := os.Getenv("FAUCET_ADDRESS")
		if faucetAddr == "" {
			faucetAddr = "cosmos1utz6mucemw5twtmjwxutja82fdeqc59awf4z8t"
		}
		
		faucetPrivKey := os.Getenv("FAUCET_PRIVATE_KEY")

		// Get faucet account
		faucetAccAddr, err := sdk.AccAddressFromBech32(faucetAddr)
		if err != nil {
			resp := faucetResp{Success: false, Message: "Invalid faucet address"}
			rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
			return
		}

		coins, err := sdk.ParseCoins(faucetAmount)
		if err != nil {
			resp := faucetResp{Success: false, Message: "Invalid faucet amount"}
			rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
			return
		}

		msg := sdk.NewMsgSend(faucetAccAddr, toAddr, coins)

		// If private key is provided, sign and broadcast
		if faucetPrivKey != "" && len(faucetPrivKey) == 64 {
			// Get account info from node
			acc, err := cliCtx.GetAccount(faucetAccAddr)
			var accNum, seq uint64
			if err == nil && acc != nil {
				accNum = acc.GetAccountNumber()
				seq = acc.GetSequence()
			}
			
			// Decode private key
			privKeyBytes, err := hex.DecodeString(faucetPrivKey)
			if err != nil {
				resp := faucetResp{Success: false, Message: fmt.Sprintf("Invalid private key: %v", err)}
				rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
				return
			}
			
			var privKey secp256k1.PrivKey
			copy(privKey[:], privKeyBytes[:32])
			pubKey := privKey.PubKey()

			// Build transaction
			txBuilder := cliCtx.NewTxBuilder()
			txBuilder.SetMsgs(msg)
			txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uos", sdk.NewInt(1000))))
			txBuilder.SetGasLimit(200000)
			txBuilder.SetMemo("faucet")
			txBuilder.SetPubKey(pubKey)
			txBuilder.SetAccountNumber(accNum)
			txBuilder.SetSequence(seq)

			// Sign
			txBuilder.SetChainID("test-chain")
			tx, err := txBuilder.GetTx().ValidateBasic()
			if err != nil {
				resp := faucetResp{Success: false, Message: fmt.Sprintf("Invalid tx: %v", err)}
				rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
				return
			}

			// Sign with the private key
			signerData := sdk.SignerData{
				ChainID:       "test-chain",
				AccountNumber: accNum,
				Sequence:      seq,
			}
			
			txBytes, err := cliCtx.SignStdTx("", faucetAddr, txBuilder.GetTx(), false)
			if err != nil {
				// Fallback: try direct broadcast
				resp := faucetResp{Success: false, Message: fmt.Sprintf("Please configure faucet with proper key: %v", err)}
				rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
				return
			}

			// Broadcast
			res, err := cliCtx.BroadcastTx(txBytes)
			if err != nil {
				resp := faucetResp{Success: false, Message: fmt.Sprintf("Broadcast failed: %v", err)}
				rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
				return
			}

			resp := faucetResp{Success: true, Message: fmt.Sprintf("Sent %s to %s", faucetAmount, req.Address), TxHash: res.TxHash}
			rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
			return
		}

		// Without private key, return info
		resp := faucetResp{Success: true, Message: fmt.Sprintf("Faucet account: %s. Configure FAUCET_PRIVATE_KEY env to enable transfers.", faucetAddr)}
		rest.PostProcessResponse(w, cdc, resp, cliCtx.Indent)
	}
}
