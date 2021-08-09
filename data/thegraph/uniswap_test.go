package thegraph

import (
	"encoding/json"
	"io"
	"reflect"
	"sort"
	"testing"

	"github.com/mahamed-belkheir/uniswrap/data"
)

func assert(got, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nexpected: %v \n got: %v", expected, got)
	}
}

func ok(err interface{}, t *testing.T) {
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

type poolQueryTestCase struct {
	firstPool  poolQueryResponse
	secondPool poolQueryResponse
	expected   []data.Pool
}

var tokens = []data.Token{
	{ID: "1", Name: "BTC"},
	{ID: "2", Name: "ETH"},
	{ID: "3", Name: "BCH"},
	{ID: "4", Name: "ADA"},
	{ID: "5", Name: "DOGE"},
}

var pools = []data.Pool{
	{ID: "1", FirstToken: tokens[0], SecondToken: tokens[1]},
	{ID: "2", FirstToken: tokens[1], SecondToken: tokens[2]},
	{ID: "3", FirstToken: tokens[3], SecondToken: tokens[4]},
	{ID: "4", FirstToken: tokens[4], SecondToken: tokens[1]},
	{ID: "5", FirstToken: tokens[0], SecondToken: tokens[2]},
}

type mockGql struct {
	response *[][]byte
}

func (m mockGql) add(v interface{}) {
	val, _ := json.Marshal(v)
	*m.response = append(*m.response, val)
}

func (m mockGql) Send(query io.Reader, v interface{}) error {
	val := (*m.response)[0]
	*m.response = (*m.response)[1:]
	return json.Unmarshal(val, v)
}

func TestPoolQuery(t *testing.T) {
	tb := []poolQueryTestCase{
		{
			firstPool: poolQueryResponse{
				Data: struct {
					Pools []data.Pool "json:\"pools\""
				}{
					Pools: pools[:2],
				},
			},
			secondPool: poolQueryResponse{
				struct {
					Pools []data.Pool "json:\"pools\""
				}{
					Pools: pools[2:],
				},
			},
			expected: pools,
		},
		{
			firstPool: poolQueryResponse{
				struct {
					Pools []data.Pool "json:\"pools\""
				}{
					Pools: pools[:1],
				},
			},
			secondPool: poolQueryResponse{
				struct {
					Pools []data.Pool "json:\"pools\""
				}{
					Pools: []data.Pool{},
				},
			},
			expected: pools[:1],
		},
	}

	mock := mockGql{&[][]byte{}}
	u := uniswap{
		mock,
	}
	for _, test := range tb {
		mock.add(test.firstPool)
		mock.add(test.secondPool)
		result, err := u.PoolsWithAsset("this does nothing")
		ok(err, t)
		assert(result, test.expected, t)
	}
}

type assetSwapVolumeTestCase struct {
	data     tokenDayQueryResponse
	expected float64
}

var dayData = []tokenDayQueryResponse{
	{
		struct {
			TokenDaysDatas []struct {
				VolumeUSD string "json:\"volumeUSD\""
			} "json:\"tokenDayDatas\""
		}{
			TokenDaysDatas: []struct {
				VolumeUSD string "json:\"volumeUSD\""
			}{
				{"100"},
				{"200"},
				{"300"},
				{"400"},
			},
		},
	},
	{
		struct {
			TokenDaysDatas []struct {
				VolumeUSD string "json:\"volumeUSD\""
			} "json:\"tokenDayDatas\""
		}{
			TokenDaysDatas: []struct {
				VolumeUSD string "json:\"volumeUSD\""
			}{
				{"500"},
				{"250"},
			},
		},
	},
	{
		struct {
			TokenDaysDatas []struct {
				VolumeUSD string "json:\"volumeUSD\""
			} "json:\"tokenDayDatas\""
		}{
			TokenDaysDatas: []struct {
				VolumeUSD string "json:\"volumeUSD\""
			}{},
		},
	},
}

func TestAssetSwapVolume(t *testing.T) {
	tb := []assetSwapVolumeTestCase{
		{
			data:     dayData[0],
			expected: 1000,
		},
		{
			data:     dayData[1],
			expected: 750,
		},
		{
			data:     dayData[2],
			expected: 0,
		},
	}
	mock := mockGql{&[][]byte{}}
	u := uniswap{
		mock,
	}
	for _, test := range tb {
		mock.add(test.data)
		result, err := u.AssetSwapVolume("this does nothing", 16000000, 16000000)
		ok(err, t)
		assert(result, test.expected, t)
	}
}

type AssetsSwappedInBlockTestCase struct {
	trx      transactionQueryResponse
	expected []data.Token
}

var swaps = []data.Swap{
	{ID: "1", AmountUSD: "500", FirstToken: tokens[0], SecondToken: tokens[1]},
	{ID: "1", AmountUSD: "500", FirstToken: tokens[2], SecondToken: tokens[1]},
	{ID: "1", AmountUSD: "500", FirstToken: tokens[0], SecondToken: tokens[2]},
	{ID: "1", AmountUSD: "500", FirstToken: tokens[3], SecondToken: tokens[2]},
	{ID: "1", AmountUSD: "500", FirstToken: tokens[4], SecondToken: tokens[3]},
	{ID: "1", AmountUSD: "500", FirstToken: tokens[4], SecondToken: tokens[2]},
}

var transactionQueries = []transactionQueryResponse{
	{Data: struct {
		Transactions []struct {
			ID          string      "json:\"id\""
			BlockNumber string      "json:\"blockNumber\""
			Swaps       []data.Swap "json:\"swaps\""
		} "json:\"transactions\""
	}{
		Transactions: []struct {
			ID          string      "json:\"id\""
			BlockNumber string      "json:\"blockNumber\""
			Swaps       []data.Swap "json:\"swaps\""
		}{
			{
				ID:          "1",
				BlockNumber: "123456",
				Swaps:       swaps,
			},
			{
				ID:          "1",
				BlockNumber: "123456",
				Swaps:       swaps[3:],
			},
		},
	}},
	{Data: struct {
		Transactions []struct {
			ID          string      "json:\"id\""
			BlockNumber string      "json:\"blockNumber\""
			Swaps       []data.Swap "json:\"swaps\""
		} "json:\"transactions\""
	}{
		Transactions: []struct {
			ID          string      "json:\"id\""
			BlockNumber string      "json:\"blockNumber\""
			Swaps       []data.Swap "json:\"swaps\""
		}{
			{
				ID:          "1",
				BlockNumber: "123456",
				Swaps:       swaps[:1],
			},
			{
				ID:          "1",
				BlockNumber: "123456",
				Swaps:       swaps[1:2],
			},
		},
	}},
}

type sortTokens []data.Token

func (a sortTokens) Len() int           { return len(a) }
func (a sortTokens) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortTokens) Less(i, j int) bool { return a[i].ID < a[j].ID }

func TestAssetsInBlock(t *testing.T) {
	tb := []AssetsSwappedInBlockTestCase{
		{
			trx:      transactionQueries[0],
			expected: tokens,
		},
		{
			trx:      transactionQueries[1],
			expected: tokens[:3],
		},
	}
	mock := mockGql{&[][]byte{}}
	u := uniswap{
		mock,
	}
	for _, test := range tb {
		mock.add(test.trx)
		var res sortTokens
		res, err := u.AssetsSwappedInBlock("12345")
		ok(err, t)
		sort.Slice(res, res.Less)
		assert([]data.Token(res), test.expected, t)
	}
}

type sortSwaps []data.Swap

func (a sortSwaps) Len() int           { return len(a) }
func (a sortSwaps) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortSwaps) Less(i, j int) bool { return a[i].ID < a[j].ID }

type SwapsInBlockTestCase struct {
	trx      transactionQueryResponse
	expected []data.Swap
}

func TestSwapsInBlock(t *testing.T) {
	tb := []SwapsInBlockTestCase{
		{transactionQueries[0], append(swaps, swaps[3:]...)},
		{transactionQueries[1], swaps[:2]},
	}
	mock := mockGql{&[][]byte{}}
	u := uniswap{
		mock,
	}
	for _, test := range tb {
		mock.add(test.trx)
		var res sortSwaps
		res, err := u.SwapsInBlock("12345")
		ok(err, t)
		sort.Slice(res, res.Less)
		assert([]data.Swap(res), test.expected, t)
	}
}
