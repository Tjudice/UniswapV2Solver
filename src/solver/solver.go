package solver

import "math/big"

// Problem:
// Given a list of tokens, and a set of cost functions to trade that token, find profitable routes

type BalanceSource interface {
	Id() string
	CurrencyId() string
	Balance() *big.Int
	Update(b *big.Int) *big.Int
}

type ArbitrageSolver struct {
	balances []BalanceSource
}

func (a *ArbitrageSolver) AddBalanceSource(b BalanceSource) {
	a.balances = append(a.balances, b)
}
