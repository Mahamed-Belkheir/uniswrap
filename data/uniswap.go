package data

type Uniswap interface {
	PoolsWithAsset(id string) ([]Pool, error)
	AssetSwapVolume(id string, timeRangeStart, timeRangeEnd int64) (float64, error)
	SwapsInBlock(id string)
	AssetsSwappedInBlock(id string)
}

type Pool struct {
	ID          string `json:"id"`
	FirstToken  Token  `json:"token0"`
	SecondToken Token  `json:"token1"`
}

type Token struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
