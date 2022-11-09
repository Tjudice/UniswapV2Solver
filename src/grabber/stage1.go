package grabber

import (
	"UniswapV2Solver/generated/uniswap"
	"UniswapV2Solver/src/data"
	"UniswapV2Solver/src/data/arango"
	"UniswapV2Solver/src/meta"
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// type Stage interface {
// 	GetLastUpdatedBlock(ctx context.Context) (int, error)
// 	RunStage(ctx context.Context, cl *ethclient.Client, o *StageOptions) error
// 	Name() string
// 	StageIndex() int
// 	Options() *RunOpts
// }

// Filter for all PairCreated events
// Stage2 Must be run after Stage1
// Must log progress so stages can pick up where left off
type Stage1 struct {
	CurrentBlock int
	db           *arango.DB
}

func (s *Stage1) Name() string {
	return "PairCreatedEvents"
}

func (s *Stage1) StageIndex() int {
	return 1
}

func (s *Stage1) GetLastUpdatedBlock(ctx context.Context) (int, error) {
	return GetLastBlockForStage(ctx, s.db, 1)
}

func (s *Stage1) GetAddresses(ctx context.Context) ([][]common.Address, error) {
	return [][]common.Address{}, nil
}

func (s *Stage1) RunStage(ctx context.Context, cl *ethclient.Client, o *StageOptions) error {
	bh := data.NewBatchHandler(ctx, s.db.PairCreatedEvent, 500)
	defer bh.Close()
	caller, err := uniswap.NewUniswapV2Factory(meta.UniswapV2FactoryAddress, cl)
	if err != nil {
		return err
	}
	pairs, err := caller.FilterPairCreated(&bind.FilterOpts{
		Start:   o.BlockStart,
		End:     &o.BlockEnd,
		Context: ctx,
	}, nil, nil)
	if err != nil {
		return err
	}
	for pairs.Next() {
		dat := &data.PairCreatedEvent{
			Token0: pairs.Event.Token0,
			Token1: pairs.Event.Token1,
			Pair:   pairs.Event.Pair,
			PairId: pairs.Event.Arg3.String(),
			EventMetaData: data.EventMetaData{
				Block:            pairs.Event.Raw.BlockNumber,
				TransactionIndex: pairs.Event.Raw.TxIndex,
				LogIndex:         pairs.Event.Raw.Index,
				Address:          pairs.Event.Raw.Address,
			},
		}
		err := bh.Write(dat)
		if err != nil {
			return err
		}
	}
	caller, err = uniswap.NewUniswapV2Factory(meta.SushiV2FactoryAddress, cl)
	if err != nil {
		return err
	}
	pairs, err = caller.FilterPairCreated(&bind.FilterOpts{
		Start:   o.BlockStart,
		End:     &o.BlockEnd,
		Context: ctx,
	}, nil, nil)
	if err != nil {
		return err
	}
	for pairs.Next() {
		dat := &data.PairCreatedEvent{
			Token0: pairs.Event.Token0,
			Token1: pairs.Event.Token1,
			Pair:   pairs.Event.Pair,
			PairId: pairs.Event.Arg3.String(),
			EventMetaData: data.EventMetaData{
				Block:            pairs.Event.Raw.BlockNumber,
				TransactionIndex: pairs.Event.Raw.TxIndex,
				LogIndex:         pairs.Event.Raw.Index,
				Address:          pairs.Event.Raw.Address,
			},
		}
		err := bh.Write(dat)
		if err != nil {
			return err
		}
	}
	return nil
}