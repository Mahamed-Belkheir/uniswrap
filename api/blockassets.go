package api

import (
	"net/http"

	"github.com/mahamed-belkheir/uniswrap/data"
)

type blockAssetsHandler struct {
	uniswap data.Uniswap
}

var _ http.Handler = blockAssetsHandler{}

func (a blockAssetsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	blockNumber, ok := checkQuery("id", w, r)
	if !ok {
		return
	}
	assets, err := a.uniswap.AssetsSwappedInBlock(blockNumber)
	if err != nil {
		m{
			"status":  "error",
			"message": err.Error(),
		}.send(400, w)
		return
	}
	m{
		"status": "success",
		"data": m{
			"assets": assets,
		},
	}.send(200, w)
}
