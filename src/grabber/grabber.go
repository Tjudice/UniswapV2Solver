package grabber

import (
	"UniswapV2Solver/generated/uniswap"
	"UniswapV2Solver/src/concurrency"
	"UniswapV2Solver/src/evts"
	"UniswapV2Solver/src/meta"
	"context"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Grabber struct {
	c              *ethclient.Client
	uniV2Factory   *uniswap.UniswapV2Factory
	sushiV2Factory *uniswap.UniswapV2Factory
	Pairs          []*Pair
}

func NewGrabber(c *ethclient.Client) (*Grabber, error) {
	var err error
	g := &Grabber{
		c: c,
	}
	g.uniV2Factory, err = uniswap.NewUniswapV2Factory(meta.UniswapV2FactoryAddress, c)
	if err != nil {
		return nil, err
	}
	g.sushiV2Factory, err = uniswap.NewUniswapV2Factory(meta.SushiV2FactoryAddress, c)
	return g, err
}

func (g *Grabber) UniFactory() *uniswap.UniswapV2Factory {
	return g.uniV2Factory
}

func (g *Grabber) SushiFactory() *uniswap.UniswapV2Factory {
	return g.sushiV2Factory
}

type SolvedSwap struct {
}

type Pair struct {
	Protocol     string
	factory      *uniswap.UniswapV2Factory
	Address      common.Address
	Token0       common.Address
	Token1       common.Address
	CreatedEvent *evts.PairCreated
}

func (g *Grabber) SolveTokenInBlockRange(ctx context.Context, token common.Address, start, end uint64) ([]SolvedSwap, error) {
	pairs, err := g.GetPairs(ctx, token, end)
	if err != nil {
		return nil, err
	}
	log.Println(pairs)
	return nil, nil
}

type blockRange struct {
	start uint64
	end   uint64
}

func (g *Grabber) GetPairs(ctx context.Context, token common.Address, end uint64) ([]*Pair, error) {
	conc := concurrency.NewConcurrentCaller[[]*Pair](10)
	uniBlocks := divideBlocks(meta.UniswapV2DeployBlock, end, 2000)
	sushiBlocks := divideBlocks(meta.SushiswapV2DeployBlock, end, 2000)
	for _, br := range uniBlocks {
		conc.AddCall(func() ([]*Pair, error) {
			filter, err := g.uniV2Factory.FilterPairCreated(&bind.FilterOpts{
				Context: ctx,
				Start:   br.start,
				End:     &br.end,
			}, []common.Address{token}, []common.Address{})
			if err != nil {
				return nil, err
			}
			ret := []*Pair{}
			for filter.Next() {
				ret = append(ret, &Pair{
					Protocol: "Uniswap",
					factory:  g.uniV2Factory,
					Address:  filter.Event.Pair,
					Token0:   filter.Event.Token0,
					Token1:   filter.Event.Token1,
					CreatedEvent: &evts.PairCreated{
						Token0: filter.Event.Token0.Hex(),
						Token1: filter.Event.Token1.Hex(),
						Pair:   filter.Event.Pair.Hex(),
						PairId: filter.Event.Arg3,
						Raw: &evts.EventMetaData{
							Block:            filter.Event.Raw.BlockNumber,
							TransactionIndex: filter.Event.Raw.TxIndex,
							LogIndex:         filter.Event.Raw.Index,
						},
					},
				})
			}
			filter2, err := g.uniV2Factory.FilterPairCreated(&bind.FilterOpts{
				Context: ctx,
				Start:   br.start,
				End:     &br.end,
			}, []common.Address{}, []common.Address{token})
			if err != nil {
				return nil, err
			}
			for filter2.Next() {
				ret = append(ret, &Pair{
					Protocol: "Uniswap",
					factory:  g.uniV2Factory,
					Address:  filter2.Event.Pair,
					Token0:   filter2.Event.Token0,
					Token1:   filter2.Event.Token1,
					CreatedEvent: &evts.PairCreated{
						Token0: filter2.Event.Token0.Hex(),
						Token1: filter2.Event.Token1.Hex(),
						Pair:   filter2.Event.Pair.Hex(),
						PairId: filter2.Event.Arg3,
						Raw: &evts.EventMetaData{
							Block:            filter2.Event.Raw.BlockNumber,
							TransactionIndex: filter2.Event.Raw.TxIndex,
							LogIndex:         filter2.Event.Raw.Index,
						},
					},
				})
			}
			return ret, nil
		})
	}
	for _, br := range sushiBlocks {
		conc.AddCall(func() ([]*Pair, error) {
			filter, err := g.uniV2Factory.FilterPairCreated(&bind.FilterOpts{
				Context: ctx,
				Start:   br.start,
				End:     &br.end,
			}, []common.Address{token}, []common.Address{})
			if err != nil {
				return nil, err
			}
			ret := []*Pair{}
			for filter.Next() {
				ret = append(ret, &Pair{
					Protocol: "SushiSwap",
					factory:  g.sushiV2Factory,
					Address:  filter.Event.Pair,
					Token0:   filter.Event.Token0,
					Token1:   filter.Event.Token1,
					CreatedEvent: &evts.PairCreated{
						Token0: filter.Event.Token0.Hex(),
						Token1: filter.Event.Token1.Hex(),
						Pair:   filter.Event.Pair.Hex(),
						PairId: filter.Event.Arg3,
						Raw: &evts.EventMetaData{
							Block:            filter.Event.Raw.BlockNumber,
							TransactionIndex: filter.Event.Raw.TxIndex,
							LogIndex:         filter.Event.Raw.Index,
						},
					},
				})
			}
			filter2, err := g.uniV2Factory.FilterPairCreated(&bind.FilterOpts{
				Context: ctx,
				Start:   br.start,
				End:     &br.end,
			}, []common.Address{}, []common.Address{token})
			if err != nil {
				return nil, err
			}
			for filter2.Next() {
				ret = append(ret, &Pair{
					Protocol: "SushiSwap",
					factory:  g.sushiV2Factory,
					Address:  filter2.Event.Pair,
					Token0:   filter2.Event.Token0,
					Token1:   filter2.Event.Token1,
					CreatedEvent: &evts.PairCreated{
						Token0: filter2.Event.Token0.Hex(),
						Token1: filter2.Event.Token1.Hex(),
						Pair:   filter2.Event.Pair.Hex(),
						PairId: filter2.Event.Arg3,
						Raw: &evts.EventMetaData{
							Block:            filter2.Event.Raw.BlockNumber,
							TransactionIndex: filter2.Event.Raw.TxIndex,
							LogIndex:         filter2.Event.Raw.Index,
						},
					},
				})
			}
			return ret, nil
		})
	}
	ret := []*Pair{}
	pairs, err := conc.Run()
	if err != nil {
		return nil, err
	}
	for _, p := range pairs {
		ret = append(ret, p...)
	}
	return ret, nil
}

func divideBlocks(start, end, sz uint64) []*blockRange {
	var blocks []*blockRange
	for i := start; i < end; i += sz {
		if i+sz > end {
			blocks = append(blocks, &blockRange{
				start: i,
				end:   end,
			})
			break
		}
		blocks = append(blocks, &blockRange{
			start: i,
			end:   i + sz - 1,
		})
	}
	return blocks
}
