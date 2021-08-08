# Uniswarp

Uniswarp is a Uniswap v3 API wrapper, it covers the following features:

- Find pools that contain a token as one of the pairs
- Find the tokens volume in USD within a time range
- Find the tokens that were swapped in a certain block
- Find the swaps that occured in a certain block

### How to run

run `ADDRESS=127.0.0.1:8000 go run cmd/main.go`
you can change address to anything you want, or omit it and run `go run cmd/main.go` to use `127.0.0.1:8000` by default

### Things possible to improve on

- Use more libraries (http frameworks, validation libraries, graphql clients)
- Add a mock API for integration testing
- Refactor the graph data provider methods to split data fetching from data processing (to do unit testing)
