package api

import (
	"log"
	"net/http"

	"github.com/mahamed-belkheir/uniswrap"
	"github.com/mahamed-belkheir/uniswrap/data"
)

func RunWebServer(config uniswrap.Config, uniswap data.Uniswap) {
	router := http.ServeMux{}
	router.Handle("/asset/pools", assetPoolsHandler{uniswap: uniswap})
	log.Printf("server listening at %s", config.Address)
	http.ListenAndServe(config.Address, &router)
}
