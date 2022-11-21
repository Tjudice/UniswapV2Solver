package grabber

import (
	"UniswapV2Solver/src/data"
	"UniswapV2Solver/src/data/arango"
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"git.tuxpa.in/a/zlog/log"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Runner struct {
	db            *arango.DB
	cl            *ethclient.Client
	stages        map[int][]Stage
	maxBlockRange int
	maxRetries    int
}

type StageOptions struct {
	BlockStart uint64
	BlockEnd   uint64
	Addresses  []common.Address
}

type Stage interface {
	Name() string
	StageIndex() int
	GetAddresses(ctx context.Context, blockMax int) ([][]common.Address, error)
	GetLastUpdatedBlock(ctx context.Context) (int, error)
	RunStage(ctx context.Context, cl *ethclient.Client, o *StageOptions) error
}

func NewStageRunner(db *arango.DB, cl *ethclient.Client, maxBlockRange int, maxRetries int) *Runner {
	return &Runner{
		db:            db,
		cl:            cl,
		stages:        make(map[int][]Stage),
		maxBlockRange: maxBlockRange,
		maxRetries:    maxRetries,
	}
}

func (r *Runner) AddStage(s ...Stage) {
	for _, st := range s {
		r.stages[st.StageIndex()] = append(r.stages[st.StageIndex()], st)
	}
}

func (r *Runner) RunStages(ctx context.Context) error {
	stageKeys := []int{}
	for stageKey := range r.stages {
		stageKeys = append(stageKeys, stageKey)
	}
	sort.SliceStable(stageKeys, func(i, j int) bool {
		return stageKeys[i] < stageKeys[j]
	})
	currBlock, err := r.cl.BlockNumber(ctx)
	if err != nil {
		return err
	}
	for i, s := range stageKeys {
		start, err := r.stages[s][0].GetLastUpdatedBlock(ctx)
		if err != nil {
			return err
		}
		log.Println(start)
		log.Info().Int("Stage Number", s).Str("Start Time", time.Now().String()).Int("Start Block", start).Msg("Starting Stage")
		for _, stage := range r.stages[s] {
			log.Info().Str("Stage Name", stage.Name()).Msg("Starting Stage")
			stageOpts := r.divideStageBlocks(start, int(currBlock))

			wg := errgroup.Group{}
			wg.SetLimit(12)
			successfulBlocks := []int{}
			mut := sync.Mutex{}

			for _, opt := range stageOpts {
				addrs, err := stage.GetAddresses(ctx, int(opt.BlockEnd))
				if err != nil {
					return err
				}
				if len(addrs) == 0 {
					addrs = [][]common.Address{{}}
				}
				for _, addr := range addrs {
					optFixed := &StageOptions{
						BlockStart: opt.BlockStart,
						BlockEnd:   opt.BlockEnd,
						Addresses:  addr,
					}
					wg.Go(func() error {
						log.Info().Str("Stage", stage.Name()).Int("Start Block", int(optFixed.BlockStart)).Int("End Block", int(optFixed.BlockEnd)).Str("Time", time.Now().String()).Msg("Starting Stage Process")
						for i := 0; i < r.maxRetries; i++ {
							err := stage.RunStage(ctx, r.cl, optFixed)
							if err != nil {
								log.Error().Err(err).Msg("Error running stage")
								continue
							}
							break
						}
						if i == r.maxRetries {
							log.Info().Str("Stage", stage.Name()).Int("Start Block", int(optFixed.BlockStart)).Int("End Block", int(optFixed.BlockEnd)).Str("Time", time.Now().String()).Msg("Failed Stage Process")

							return fmt.Errorf("%d", optFixed.BlockStart)
						}
						mut.Lock()
						log.Info().Str("Stage", stage.Name()).Int("Start Block", int(optFixed.BlockStart)).Int("End Block", int(optFixed.BlockEnd)).Str("Time", time.Now().String()).Msg("Finished Stage Process")
						successfulBlocks = append(successfulBlocks, int(optFixed.BlockEnd))
						mut.Unlock()
						return nil
					})
				}
			}

			err = wg.Wait()
			sort.SliceStable(successfulBlocks, func(i, j int) bool {
				return successfulBlocks[i] < successfulBlocks[j]
			})
			if err != nil {
				maxBlockSuccess, _ := strconv.Atoi(err.Error())
				bestblock := -1
				lastChecked := -1
				for _, block := range successfulBlocks {
					if block > maxBlockSuccess {
						bestblock = block
						break
					}
					lastChecked = block
				}
				if bestblock == -1 {
					_, err := r.db.StageProgress.C().CreateDocument(ctx, &data.StageProgress{
						Stage: stage.StageIndex(),
						Block: lastChecked,
						Time:  time.Now(),
					})
					if err != nil {
						return err
					}
				}
			} else {
				_, err := r.db.StageProgress.C().CreateDocument(ctx, &data.StageProgress{
					Stage: stage.StageIndex(),
					Block: int(currBlock),
					Time:  time.Now(),
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *Runner) divideStageBlocks(start, end int) []*StageOptions {
	var stageOpts []*StageOptions

	for start < end {
		if end-start < r.maxBlockRange {
			stageOpts = append(stageOpts, &StageOptions{
				BlockStart: uint64(start),
				BlockEnd:   uint64(end),
			})
			return stageOpts
		}
		stageOpts = append(stageOpts, &StageOptions{
			BlockStart: uint64(start),
			BlockEnd:   uint64(start + r.maxBlockRange),
		})
		start += r.maxBlockRange
	}
	return stageOpts
}
