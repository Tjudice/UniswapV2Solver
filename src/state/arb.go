package state

import "math/big"

type Route struct {
	Amount0 *big.Int
	Amount1 *big.Int
	Fees    *big.Int
	Next    *Route
}

type Arb struct {
	Block            int
	TransactionIndex int
	LogIndex         int
	Profit           *big.Int
	Route            *Route
}
