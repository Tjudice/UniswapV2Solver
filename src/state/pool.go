package state

import (
	"fmt"
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

func (p *Pool) Swap(amt *TokenAmount) (*TokenAmount, error) {
	if amt.Amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("insufficient input amount")
	}
	if amt.Token.Address == p.Token0.Address {
		ret := p.zeroForOne(amt.Amount)
		p.Reserve0.Add(p.Reserve0, amt.Amount)
		p.Reserve1.Sub(p.Reserve1, ret)
		return &TokenAmount{
			Token:  p.Token1,
			Amount: ret,
		}, nil
	}
	if amt.Token.Address == p.Token1.Address {
		ret := p.oneForZero(amt.Amount)
		p.Reserve1.Add(p.Reserve1, amt.Amount)
		p.Reserve0.Sub(p.Reserve0, ret)
		return &TokenAmount{
			Token:  p.Token0,
			Amount: ret,
		}, nil
	}
	return nil, fmt.Errorf("token not in pool")
}

func (p *Pool) TestSwap(amt *TokenAmount) (*TokenAmount, error) {
	if amt.Amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("insufficient input amount")
	}
	if amt.Token.Address == p.Token0.Address {
		return &TokenAmount{
			Token:  p.Token1,
			Amount: p.zeroForOne(amt.Amount),
		}, nil
	}
	if amt.Token.Address == p.Token1.Address {
		return &TokenAmount{
			Token:  p.Token0,
			Amount: p.oneForZero(amt.Amount),
		}, nil
	}
	return nil, fmt.Errorf("token not in pool")
}

func (p *Pool) zeroForOne(token0Amount *big.Int) *big.Int {
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

func (p *Pool) oneForZero(token1Amount *big.Int) *big.Int {
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

// price = reserve1Next/reserve0Next
// price = (reserve1 + token1Out)/(reserve0 + token0In)
// price * (reserve0 + token0In) = reserve1 + token1Out
// price * reserve0 + price * token0In = reserve1 + token1Out
// token0In = (reserve1 + token1Out - price * reserve0) / price
// ((reserve0 + token0In) * 1000 - 3 * token0In) * ((reserve1 + token1Out) * 1000 - 3 * token0Out) = k
// b = (reserve1 + token1Out)/price
// c = (reserve1 + token1Out - price * reserve0) / price
func (p *Pool) MaxToken0ToPrice(price *big.Int) (in *TokenAmount, out *TokenAmount) {
	oneThousand := big.NewInt(1000)
	k := new(big.Int).Mul(oneThousand, oneThousand)
	k.Mul(k, big.NewInt(0).Mul(p.Reserve0, p.Reserve1))
	return nil, nil
}

// price = reserve0Next/reserve1Next
// price = (reserve0 + token0Out)/(reserve1 + token1In)
// price * (reserve1 + token1In) = reserve0 + token0Out
// price * reserve1 + price * token1In = reserve0 + token0Out
// token1In = (reserve0 + token0Out - price * reserve1) / price
func (p *Pool) MaxToken1ToPrice(price *big.Int) (in *TokenAmount, out *TokenAmount) {
	return nil, nil
}

func (p *Pool) Price0() *big.Int {
	if p.Reserve0.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0)
	}
	return new(big.Int).Div(p.Reserve1, p.Reserve0)
}

func (p *Pool) Price1() *big.Int {
	if p.Reserve1.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0)
	}
	return new(big.Int).Div(p.Reserve0, p.Reserve1)
}

func (p *Pool) PriceForToken(token common.Address) *big.Int {
	if token == p.Token0.Address {
		return p.Price0()
	}
	if token == p.Token1.Address {
		return p.Price1()
	}
	return big.NewInt(0)
}

func (p *Pool) Copy() *Pool {
	return &Pool{
		PairId:           big.NewInt(0).Set(p.PairId),
		ContractAddress:  p.ContractAddress,
		Token0:           p.Token0,
		Token1:           p.Token1,
		Reserve0:         big.NewInt(0).Set(p.Reserve0),
		Reserve1:         big.NewInt(0).Set(p.Reserve1),
		K:                big.NewInt(0).Set(p.K),
		LastUpdatedBlock: p.LastUpdatedBlock,
	}
}
