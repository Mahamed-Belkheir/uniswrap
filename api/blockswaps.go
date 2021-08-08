package api

import (
	"net/http"

	"github.com/mahamed-belkheir/uniswrap/data"
)

type blockSwapsHandler struct {
	uniswap data.Uniswap
}

var _ http.Handler = blockSwapsHandler{}

func (a blockSwapsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	blockNumber, ok := checkQuery("id", w, r)
	if !ok {
		return
	}
	swaps, err := a.uniswap.SwapsInBlock(blockNumber)
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
			"swaps": swaps,
		},
	}.send(200, w)
}
