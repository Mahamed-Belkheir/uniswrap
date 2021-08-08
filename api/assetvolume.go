package api

import (
	"net/http"
	"strconv"

	"github.com/mahamed-belkheir/uniswrap/data"
)

type assetVolumeHandler struct {
	uniswap data.Uniswap
}

var _ http.Handler = assetVolumeHandler{}

func (a assetVolumeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	assetId, ok := checkQuery("id", w, r)
	if !ok {
		return
	}
	start, ok := checkIntQuery("start", w, r)
	if !ok {
		return
	}
	end, ok := checkIntQuery("end", w, r)
	if !ok {
		return
	}
	volume, err := a.uniswap.AssetSwapVolume(assetId, start, end)
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
			"totalVolumeUSD": volume,
		},
	}.send(200, w)
}

func checkIntQuery(query string, w http.ResponseWriter, r *http.Request) (int64, bool) {
	i, err := strconv.ParseInt(r.URL.Query().Get(query), 10, 64)
	if err != nil {
		m{
			"status":  "error",
			"message": "failed to parse query parameter " + query + " as int: " + err.Error(),
		}.send(400, w)
		return 0, false
	}
	return i, true
}
