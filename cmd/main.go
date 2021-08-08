package main

import (
	"github.com/mahamed-belkheir/uniswrap"
	"github.com/mahamed-belkheir/uniswrap/api"
	"github.com/mahamed-belkheir/uniswrap/data/thegraph"
)

func main() {
	source := thegraph.NewUniswap()
	config := uniswrap.GetConfig()
	api.RunWebServer(config, source)
}
