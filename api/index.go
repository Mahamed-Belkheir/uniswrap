package api

import (
	"net/http"

	"github.com/mahamed-belkheir/uniswrap"
)

func RunWebServer(config uniswrap.Config) {
	router := http.ServeMux{}

	http.ListenAndServe(config.Address, &router)
}
