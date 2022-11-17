package grabber

import (
	"UniswapV2Solver/src/data"
	"UniswapV2Solver/src/data/arango"
	"UniswapV2Solver/src/evts"
	"UniswapV2Solver/src/meta"
	"context"
	"fmt"

	builder "gfx.cafe/open/arango"
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
	blk, err := GetLastBlockForStage(ctx, s.db, 2)
	if err != nil {
		blk = int(meta.UniswapV2DeployBlock)
	}
	return blk, nil
}

const PairAddressesQuery = `
	for doc in PairCreatedEvent
	filter doc.block <= %d
	let ret = doc.pair	
`

func (s *Stage2) GetAddresses(ctx context.Context, blockMax int) ([][]common.Address, error) {
	addrsString, err := builder.NewBuilder[string](s.db.D()).
		Raw(fmt.Sprintf(PairAddressesQuery, blockMax)).
		Return(ctx, "ret")
	if err != nil {
		return nil, err
	}
	addrs := make([]common.Address, len(addrsString))
	for i, a := range addrsString {
		addrs[i] = common.HexToAddress(a)
	}
	numDivisions := 50
	addrSize := len(addrs) / numDivisions
	ret := make([][]common.Address, numDivisions)
	for i := 0; i < numDivisions; i++ {
		if i != numDivisions-1 {
			ret[i] = addrs[i*addrSize : (i+1)*addrSize]
		}
		ret[i] = addrs[i*addrSize:]
	}
	return ret, nil
}

func (s *Stage2) RunStage(ctx context.Context, cl *ethclient.Client, o *StageOptions) error {
	bh := data.NewBatchHandler(ctx, s.db.SyncEvent, 100)
	defer bh.Close()
	filterer, err := evts.NewUniswapContractGroupFilterer(cl)
	if err != nil {
		return err
	}
	fo := &bind.FilterOpts{
		Start:   o.BlockStart,
		End:     &o.BlockEnd,
		Context: ctx,
	}

	evs, err := filterer.FilterSync(fo, o.Addresses)
	if err != nil {
		return err
	}
	for _, ev := range evs {
		dat := &data.SyncEvent{
			Reserve0:         ev.Reserve0.String(),
			Reserve1:         ev.Reserve1.String(),
			Block:            ev.Raw.Block,
			Transaction:      ev.Raw.Transaction,
			TransactionIndex: ev.Raw.TransactionIndex,
			LogIndex:         ev.Raw.LogIndex,
			Address:          ev.Raw.Address,
			BlockTimestamp:   ev.Raw.BlockTimestamp,
		}
		err := bh.Write(dat)
		if err != nil {
			return err
		}
	}
	return nil
}
