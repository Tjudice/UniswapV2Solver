package grabber

import (
	"UniswapV2Solver/src/data/arango"
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Runner struct {
	db     *arango.DB
	cl     *ethclient.Client
	stages map[int][]Stage
}

type RunOpts struct {
	BlockDivisionSize int
	NumConcurrent     int
}

func NewStageRunner(db *arango.DB, cl *ethclient.Client) *Runner {
	return &Runner{
		db:     db,
		cl:     cl,
		stages: make(map[int][]Stage),
	}
}

func (r *Runner) AddStage(s ...Stage) {
	for _, st := range s {
		r.stages[st.StageIndex()] = append(r.stages[st.StageIndex()], st)
	}
}

func (r *Runner) RunStages(ctx context.Context) error {
	return nil
}

type StageOptions struct {
	BlockStart uint64
	BlockEnd   uint64
	Addresses  []common.Address
}

type Stage interface {
	GetLastUpdatedBlock(ctx context.Context) (int, error)
	RunStage(ctx context.Context, cl *ethclient.Client, s *StageOptions) error
	Name() string
	StageIndex() int
	Options() *RunOpts
}
