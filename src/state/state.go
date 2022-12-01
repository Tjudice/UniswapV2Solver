package state

import (
	"UniswapV2Solver/src/evts"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
)

type UniswapV2State struct {
	// Clients
	cl *ethclient.Client

	// Data
	Pools  map[common.Address]*Pool
	Tokens map[common.Address]*Token

	// Current state index
	Block            int
	TransactionIndex int
	LogIndex         int

	// Counters
	EventCount int

	// Mutex
	mut sync.Mutex
}

func NewState(cl *ethclient.Client) *UniswapV2State {
	return &UniswapV2State{
		cl:     cl,
		Pools:  map[common.Address]*Pool{},
		Tokens: map[common.Address]*Token{},
		mut:    sync.Mutex{},
	}
}

func (s *UniswapV2State) AddTokens(tokens []data.Token) error {
	return fmt.Errorf("not implemented")
}

func (s *UniswapV2State) AddAllPairs(evt []evts.UniV2Event) error {
	wg := errgroup.Group{}
	wg.SetLimit(12)
	for _, evt := range evt {
		v, ok := evt.(*evts.PairCreated)
		if !ok {
			continue
		}
		wg.Go(func() error {
			return s.pairCreated(v)
		})
	}
	err := wg.Wait()
	return err
}

func (s *UniswapV2State) Update(evt evts.UniV2Event) error {
	switch v := evt.(type) {
	case *evts.PairCreated:
		err := s.pairCreated(v)
		if err != nil {
			return err
		}
	case *evts.Sync:
		pool, err := s.GetPool(common.HexToAddress(v.Raw.Address))
		if err != nil {
			return err
		}
		pool.Update(v.Reserve0, v.Reserve1, int(v.Raw.Block))
	default:
		return fmt.Errorf("invalid event type")
	}
	s.Block, s.TransactionIndex, s.LogIndex = evt.EventIndex()
	s.EventCount = s.EventCount + 1
	return nil
}

func (s *UniswapV2State) pairCreated(v *evts.PairCreated) error {
	pair := common.HexToAddress(v.Pair)
	s.mut.Lock()
	if _, ok := s.Pools[pair]; ok {
		s.mut.Unlock()
		return nil
	}
	s.mut.Unlock()
	token0addr := common.HexToAddress(v.Token0)
	token1addr := common.HexToAddress(v.Token1)
	token0, err := s.GetToken(token0addr)
	if err != nil {
		return err
	}
	token1, err := s.GetToken(token1addr)
	if err != nil {
		return err
	}
	s.mut.Lock()
	defer s.mut.Unlock()
	s.Pools[pair] = &Pool{
		PairId:           v.PairId,
		ContractAddress:  pair,
		Reserve0:         big.NewInt(0),
		Reserve1:         big.NewInt(0),
		Token0:           token0,
		Token1:           token1,
		LastUpdatedBlock: int(v.Raw.Block),
	}
	return nil
}

func (s *UniswapV2State) GetToken(address common.Address) (*Token, error) {
	s.mut.Lock()
	defer s.mut.Unlock()
	if token, ok := s.Tokens[address]; ok {
		return token, nil
	}
	// caller, _ := erc20.NewErc20Caller(address, s.cl)
	// dec, err := caller.Decimals(nil)
	// if err != nil {
	// 	return nil, err
	// }
	// s.Tokens[address] = NewToken(address, int(dec))
	s.Tokens[address] = NewToken(address, 18)
	return s.Tokens[address], nil
}

func (s *UniswapV2State) GetPool(address common.Address) (*Pool, error) {
	s.mut.Lock()
	defer s.mut.Unlock()
	if pool, ok := s.Pools[address]; ok {
		return pool, nil
	}
	return nil, fmt.Errorf("pool not found")
}

func (s *UniswapV2State) PrintStateSummary() {
	fmt.Printf("============= State ==============\n")
	fmt.Printf("Block: %d\n", s.Block)
	fmt.Printf("Transaction Index: %d\n", s.TransactionIndex)
	fmt.Printf("Log Index: %d\n", s.LogIndex)
	fmt.Printf("Pools: %d\n", len(s.Pools))
	fmt.Printf("Tokens: %d\n", len(s.Tokens))
	fmt.Printf("Event Count: %d\n", s.EventCount)
}
