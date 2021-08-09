package data

//Uniswap is an interface for fetching data from uniswap
type Uniswap interface {
	PoolsWithAsset(id string) ([]Pool, error)
	AssetSwapVolume(id string, timeRangeStart, timeRangeEnd int64) (float64, error)
	SwapsInBlock(id string) ([]Swap, error)
	AssetsSwappedInBlock(id string) ([]Token, error)
}

//Pool represent a uniswap pool of two tokens
type Pool struct {
	ID          string `json:"id"`
	FirstToken  Token  `json:"token0"`
	SecondToken Token  `json:"token1"`
}

//Token is a token on the uniswap exchange
type Token struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//Swap is a swap operation between two tokens in a token pool
type Swap struct {
	AmountUSD   string `json:"amountUSD"`
	ID          string `json:"id"`
	FirstToken  Token  `json:"token0"`
	SecondToken Token  `json:"token1"`
}
