package api

import (
	"net/http"

	"github.com/mahamed-belkheir/uniswrap/data"
)

type assetPoolsHandler struct {
	uniswap data.Uniswap
}

var _ http.Handler = assetPoolsHandler{}

func (a assetPoolsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	assetId, ok := checkQuery("id", w, r)
	if !ok {
		return
	}
	pools, err := a.uniswap.PoolsWithAsset(assetId)
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
			"pools": pools,
		},
	}.send(200, w)

}
