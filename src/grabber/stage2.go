package grabber

import (
	"UniswapV2Solver/generated/uniswap"
	"UniswapV2Solver/src/data"
	"UniswapV2Solver/src/data/arango"
	"UniswapV2Solver/src/meta"
	"context"

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
// Stage2 Must be run after Stage2
// Must log progress so stages can pick up where left off
type Stage2 struct {
	db *arango.DB
}

func NewStage2(db *arango.DB) *Stage2 {
	return &Stage2{db: db}
}

func (s *Stage2) Name() string {
	return "PoolEvents"
}

func (s *Stage2) StageIndex() int {
	return 2
}

func (s *Stage2) GetLastUpdatedBlock(ctx context.Context) (int, error) {
	blk, err := GetLastBlockForStage(ctx, s.db, 1)
	if err != nil {
		blk = int(meta.UniswapV2DeployBlock)
	}
	return blk, nil
}

func (s *Stage2) GetAddresses(ctx context.Context) ([][]common.Address, error) {
	return [][]common.Address{}, nil
}

func (s *Stage2) RunStage(ctx context.Context, cl *ethclient.Client, o *StageOptions) error {
	bh := data.NewBatchHandler(ctx, s.db.PairCreatedEvent, 1000)
	defer bh.Close()
	caller, err := uniswap.NewUniswapV2Pair(meta.UniswapV2FactoryAddress, cl)
	if err != nil {
		return err
	}
	return nil
}
