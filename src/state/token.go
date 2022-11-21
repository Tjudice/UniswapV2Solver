package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	Address  common.Address
	Pools    []*Pool
	Decimals int
}

func NewToken(addr common.Address, decimals int) *Token {
	return &Token{
		Address:  addr,
		Decimals: decimals,
		Pools:    []*Pool{},
	}
}

type TokenAmount struct {
	Token  *Token
	Amount *big.Int
}

func (t *Token) Update() {

}
