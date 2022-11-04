package evts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// event PairCreated(address indexed token0, address indexed token1, address pair, uint);

type PairCreated struct {
	Block            uint64
	TransactionIndex uint
	LogIndex         uint

	Token0 common.Address
	Token1 common.Address
	Pair   common.Address
	PairId *big.Int
}
