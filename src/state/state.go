package state

import (
	"UniswapV2Solver/generated/erc20"
	"UniswapV2Solver/src/evts"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type UniswapV2State struct {
	// Clients
	cl *ethclient.Client

	// Data
	Pools  map[common.Address]*Pool
	Tokens map[common.Address]*Token

	Block            int
	TransactionIndex int
	LogIndex         int
}

func NewState(cl *ethclient.Client) *UniswapV2State {
	return &UniswapV2State{
		cl:     cl,
		Pools:  map[common.Address]*Pool{},
		Tokens: map[common.Address]*Token{},
	}
}

func (s *UniswapV2State) Update(evt evts.UniV2Event) error {
	switch v := evt.(type) {
	case *evts.PairCreated:
		pair := common.HexToAddress(v.Pair)
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
		s.Pools[pair] = &Pool{
			PairId:           v.PairId,
			ContractAddress:  pair,
			Reserve0:         big.NewInt(0),
			Reserve1:         big.NewInt(0),
			Token0:           token0,
			Token1:           token1,
			LastUpdatedBlock: int(v.Raw.Block),
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
	return nil
}

func (s *UniswapV2State) GetToken(address common.Address) (*Token, error) {
	if token, ok := s.Tokens[address]; ok {
		return token, nil
	}
	caller, _ := erc20.NewErc20Caller(address, s.cl)
	dec, err := caller.Decimals(nil)
	if err != nil {
		return nil, err
	}
	s.Tokens[address] = NewToken(address, int(dec))
	return s.Tokens[address], nil
}

func (s *UniswapV2State) GetPool(address common.Address) (*Pool, error) {
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
}
