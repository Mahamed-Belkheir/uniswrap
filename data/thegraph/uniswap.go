package thegraph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/mahamed-belkheir/uniswrap/data"
)

type uniswap struct {
	url    string
	client http.Client
}

func NewUniswap() uniswap {
	return uniswap{
		url:    "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3",
		client: http.Client{},
	}
}

var _ data.Uniswap = uniswap{}

func (u uniswap) PoolsWithAsset(id string) ([]data.Pool, error) {
	tok0, err := u.poolQuery(id, "token0")
	if err != nil {
		return nil, err
	}
	tok1, err := u.poolQuery(id, "token1")
	if err != nil {
		return nil, err
	}
	return append(tok0, tok1...), nil
}

type poolQueryResponse struct {
	Data struct {
		Pools []data.Pool `json:"pools"`
	} `json:"data"`
}

func (u uniswap) poolQuery(id, key string) ([]data.Pool, error) {
	res, err := u.client.Post(u.url, "application/json", m{
		"query": fmt.Sprintf(`query {
			pools(where: { %s: "%s"}) {
				id,
				token0 { id, name },
				token1 { id, name }
			}
		}
		`, key, id),
	}.json())
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	var result poolQueryResponse
	err = json.Unmarshal(rawBody, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json response: %w", err)
	}
	return result.Data.Pools, nil
}

type tokenDayQueryResponse struct {
	Data struct {
		TokenDaysDatas []struct {
			VolumeUSD string `json:"volumeUSD"`
		} `json:"tokenDayDatas"`
	} `json:"data"`
}

func (u uniswap) AssetSwapVolume(id string, start int64, end int64) (float64, error) {
	res, err := u.client.Post(u.url, "application/json", m{
		"query": fmt.Sprintf(`query {
			tokenDayDatas(where: { token: "%s", volume_gt: 0, date_gte: %v, date_lte: %v}) {
				volumeUSD
			  }
		}
		`, id, start, end),
	}.json())
	if err != nil {
		return 0, fmt.Errorf("error sending request: %w", err)
	}
	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response: %w", err)
	}
	var result tokenDayQueryResponse
	err = json.Unmarshal(rawBody, &result)
	if err != nil {
		return 0, fmt.Errorf("error unmarshaling json response: %w", err)
	}

	var total float64 = 0
	for _, res := range result.Data.TokenDaysDatas {
		num, err := strconv.ParseFloat(res.VolumeUSD, 64)
		if err != nil {
			return 0, fmt.Errorf("error casting volume to float: \n %s", res)
		}
		total += num
	}
	return total, nil
}

func (u uniswap) SwapsInBlock(id string) {

}

func (u uniswap) AssetsSwappedInBlock(id string) {

}
