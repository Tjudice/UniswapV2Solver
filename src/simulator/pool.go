package simulator

import (
	"UniswapV2Solver/src/evts"
	"UniswapV2Solver/src/meta"
	"UniswapV2Solver/src/simulator/UQ112x112"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const MINIMUM_LIQUIDITY int64 = 1000

type Pool struct {
	PairId int `json:"pair_id"`

	Token0 common.Address `json:"token0"`
	Token1 common.Address `json:"token1"`

	PoolAddress common.Address `json:"pool_address"`

	TotalSupply *big.Int `json:"total_supply"`

	Reserve0 *big.Int `json:"reserve0"`
	Reserve1 *big.Int `json:"reserve1"`

	Balance0 *big.Int `json:"balance0"`
	Balance1 *big.Int `json:"balance1"`

	Price0CumulativeLast *big.Int `json:"price0_cumulative_last"`
	Price1CumulativeLast *big.Int `json:"price1_cumulative_last"`
	KLast                *big.Int `json:"k_last"`

	BlockTimestampLast time.Time `json:"block_timestamp_last"`

	CreateEvent *evts.PairCreated `json:"create_event"`

	Positions *Position `json:"positions"`
}

func NewPool(evt *evts.PairCreated) *Pool {
	return &Pool{
		PairId:               int(evt.PairId.Int64()),
		Token0:               common.HexToAddress(evt.Token0),
		Token1:               common.HexToAddress(evt.Token1),
		PoolAddress:          common.HexToAddress(evt.Pair),
		Reserve0:             big.NewInt(0),
		Reserve1:             big.NewInt(0),
		Balance0:             big.NewInt(0),
		Balance1:             big.NewInt(0),
		Price0CumulativeLast: big.NewInt(0),
		Price1CumulativeLast: big.NewInt(0),
		KLast:                big.NewInt(0),
		BlockTimestampLast:   evt.Raw.BlockTimestamp,
		CreateEvent:          evt,
		Positions:            &Position{},
	}
}

func (p *Pool) SimulateEvent(evt evts.UniV2Event) error {
	switch e := evt.(type) {
	case *evts.Mint:
		return p.Mint(e)
	case *evts.Burn:
		return p.Burn(e)
	case *evts.Swap:
		return p.Swap(e)
	case *evts.Sync:
		return p.Sync(e)
	default:
		return fmt.Errorf("unknown event type: %T", evt)
	}
}

func (p *Pool) Mint(e *evts.Mint) error {
	var liq *big.Int
	if p.TotalSupply.Sign() == 0 {
		liq := big.NewInt(0).Sqrt(big.NewInt(0).Mul(e.Amount0, e.Amount1))
		liq.Sub(liq, big.NewInt(MINIMUM_LIQUIDITY))
		p.Positions.Mint(meta.ZeroAddress, big.NewInt(MINIMUM_LIQUIDITY))
	} else {
		liq1 := big.NewInt(0).Div(big.NewInt(0).Mul(e.Amount0, p.Positions.TotalSupply), p.Reserve0)
		liq2 := big.NewInt(0).Div(big.NewInt(0).Mul(e.Amount1, p.Positions.TotalSupply), p.Reserve1)
		if liq1.Cmp(liq2) >= 0 {
			liq = liq2
		} else {
			liq = liq1
		}
	}
	p.Positions.Mint(common.HexToAddress(e.Sender), liq)
	return nil
}

func (p *Pool) Burn(e *evts.Burn) error {
	liquidity := big.NewInt(0).Div(big.NewInt(0).Mul(e.Amount0, p.TotalSupply), p.Balance0)
	p.Positions.Burn(e.Sender, liquidity)
	p.TotalSupply.Sub(p.TotalSupply, liquidity)
	return nil
}

func (p *Pool) Swap(e *evts.Swap) error {
	// Do I even need to do anything?
	return nil
}

func (p *Pool) Sync(e *evts.Sync) error {
	timeElapsed := e.Raw.BlockTimestamp.Unix() - p.BlockTimestampLast.Unix()
	if timeElapsed > 0 && p.Reserve0.Cmp(big.NewInt(0)) != 0 && p.Reserve1.Cmp(big.NewInt(0)) != 0 {
		p.Price0CumulativeLast.Add(p.Price0CumulativeLast, big.NewInt(0).Mul(UQ112x112.Uqdiv(UQ112x112.Encode(p.Reserve1), p.Reserve0), big.NewInt(int64(timeElapsed))))
		p.Price1CumulativeLast.Add(p.Price1CumulativeLast, big.NewInt(0).Mul(UQ112x112.Uqdiv(UQ112x112.Encode(p.Reserve0), p.Reserve1), big.NewInt(int64(timeElapsed))))
	}
	p.Reserve0.Set(e.Reserve0)
	p.Reserve1.Set(e.Reserve1)
	p.Balance0.Set(e.Reserve0)
	p.Balance1.Set(e.Reserve1)

	p.BlockTimestampLast = e.Raw.BlockTimestamp
	return nil
}
