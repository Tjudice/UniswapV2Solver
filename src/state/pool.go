package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Pool struct {
	PairId           *big.Int
	ContractAddress  common.Address
	Token0           *Token
	Token1           *Token
	Reserve0         *big.Int
	Reserve1         *big.Int
	K                *big.Int
	LastUpdatedBlock int
}

func (p *Pool) Update(reserve0, reserve1 *big.Int, block int) {
	p.Reserve0 = reserve0
	p.Reserve1 = reserve1
	p.LastUpdatedBlock = block
	p.K = new(big.Int).Mul(reserve0, reserve1)
	p.Token0.Update()
	p.Token1.Update()
}

func (p *Pool) SwapZeroForOneOutput(token0Amount *big.Int) *big.Int {
	oneThousand := big.NewInt(1000)
	k := new(big.Int).Mul(oneThousand, oneThousand)
	k.Mul(k, big.NewInt(0).Mul(p.Reserve0, p.Reserve1))
	z := new(big.Int).Add(p.Reserve0, token0Amount)
	z.Mul(z, oneThousand)
	z.Sub(z, new(big.Int).Mul(big.NewInt(3), token0Amount))
	numerator := new(big.Int).Div(k, z)
	numerator.Sub(numerator, new(big.Int).Mul(oneThousand, p.Reserve1))
	return numerator.Div(numerator, big.NewInt(997))
}

func (p *Pool) OneForZeroOutput(token1Amount *big.Int) *big.Int {
	oneThousand := big.NewInt(1000)
	k := new(big.Int).Mul(oneThousand, oneThousand)
	k.Mul(k, big.NewInt(0).Mul(p.Reserve0, p.Reserve1))
	z := new(big.Int).Add(p.Reserve1, token1Amount)
	z.Mul(z, oneThousand)
	z.Sub(z, new(big.Int).Mul(big.NewInt(3), token1Amount))
	numerator := new(big.Int).Div(k, z)
	numerator.Sub(numerator, new(big.Int).Mul(oneThousand, p.Reserve0))
	return numerator.Div(numerator, big.NewInt(997))
}
