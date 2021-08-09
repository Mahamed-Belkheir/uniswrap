package thegraph

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mahamed-belkheir/uniswrap/data"
)

type uniswap struct {
	gql gqlClient
}

/*NewUniswap provides a new instance of TheGraph API implementation for the uniswap data provider interface*/
func NewUniswap() uniswap {
	return uniswap{
		gql: httpPostGQL{
			url: "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3",
			c:   http.Client{},
		},
	}
}

var _ data.Uniswap = uniswap{}

/*PoolsWithAsset fetches all the pools that contains the token id
it makes two queries to the graph because there's no orWhere option,
so we check the first token and second tokens for the wanted token.
*/
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
	var result poolQueryResponse
	err := u.gql.Send(m{
		"query": fmt.Sprintf(`query {
			pools(where: { %s: "%s"}) {
				id,
				token0 { id, name },
				token1 { id, name }
			}
		}
		`, key, id),
	}.json(), &result)
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

/*AssetSwapVolume calculates the total volume traded within the provided date range
it fetches token day data and calculates the sum*/
func (u uniswap) AssetSwapVolume(id string, start int64, end int64) (float64, error) {
	var result tokenDayQueryResponse
	err := u.gql.Send(m{
		// notes: volume_gt: 0 so as to avoid loading token day data with no volume traded
		"query": fmt.Sprintf(`query { 
			tokenDayDatas(where: { token: "%s", volume_gt: 0, date_gte: %v, date_lte: %v}) {
				volumeUSD
			  }
		}
		`, id, start, end),
	}.json(), &result)
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

type transactionQueryResponse struct {
	Data struct {
		Transactions []struct {
			ID          string      `json:"id"`
			BlockNumber string      `json:"blockNumber"`
			Swaps       []data.Swap `json:"swaps"`
		} `json:"transactions"`
	} `json:"data"`
}

/*SwapsInBlock fetches all the swaps in the block by finding all the transactions in the block
and concating the transaction swaps*/
func (u uniswap) SwapsInBlock(id string) ([]data.Swap, error) {
	trx, err := u.queryTransactions(id)
	if err != nil {
		return nil, err
	}
	swaps := make([]data.Swap, 0)
	for _, t := range trx.Data.Transactions {
		swaps = append(swaps, t.Swaps...)
	}
	return swaps, nil
}

/*AssetsSwappedInBlock fetches the assets that were traded in the provided block number
the returned slice is a valid set without any repeated values*/
func (u uniswap) AssetsSwappedInBlock(id string) ([]data.Token, error) {
	trx, err := u.queryTransactions(id)
	if err != nil {
		return nil, err
	}
	// note: using a map as a set to get rid of repeated token entries
	tokenSet := map[string]data.Token{}
	for _, t := range trx.Data.Transactions {
		for _, s := range t.Swaps {
			tokenSet[s.FirstToken.ID] = s.FirstToken
			tokenSet[s.SecondToken.ID] = s.SecondToken
		}
	}
	tokens := make([]data.Token, len(tokenSet))
	i := 0
	for _, token := range tokenSet {
		tokens[i] = token
		i++
	}
	return tokens, nil
}

func (u uniswap) queryTransactions(id string) (transactionQueryResponse, error) {
	var trx transactionQueryResponse
	err := u.gql.Send(m{
		"query": fmt.Sprintf(`query {
			transactions(where: { blockNumber: "%v"}) {
				id,
				blockNumber,
				swaps {
				  id,
				  token0 {
					id,
					name  
				  },
				  token1 {
					id,
					name
				  },
				  amountUSD
				}
			  }
		}`, id),
	}.json(), &trx)
	if err != nil {
		return trx, fmt.Errorf("error unmarshaling json response: %w", err)
	}
	return trx, nil
}
