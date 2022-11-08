package simulator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Position struct {
	TotalSupply *big.Int `json:"total_supply"`

	BalanceOf map[common.Address]*big.Int `json:"balance_of"`
}

func (p *Position) Mint(to common.Address, amount *big.Int) {
	p.TotalSupply.Add(p.TotalSupply, amount)
	pos, ok := p.BalanceOf[to]
	if !ok {
		p.BalanceOf[to] = big.NewInt(0)
	}
	pos.Add(p.BalanceOf[to], amount)
	p.BalanceOf[to] = pos
}

func (p *Position) Burn(from common.Address, amount *big.Int) {
	if _, ok := p.BalanceOf[from]; !ok {
		return
	}
	p.TotalSupply.Sub(p.TotalSupply, amount)
	p.BalanceOf[from].Sub(p.BalanceOf[from], amount)
	if p.BalanceOf[from].Cmp(big.NewInt(0)) == 0 {
		delete(p.BalanceOf, from)
	}
}

func (p *Position) Get(addr common.Address) (*big.Int, bool) {
	balance, ok := p.BalanceOf[addr]
	return balance, ok
}
