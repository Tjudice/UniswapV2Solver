package evts

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type EventMetaData struct {
	Block            uint64
	TransactionIndex uint
	LogIndex         uint
	Address          common.Address
	BlockTimestamp   time.Time
}

// event PairCreated(address indexed token0, address indexed token1, address pair, uint);
type PairCreated struct {
	Token0 common.Address
	Token1 common.Address
	Pair   common.Address
	PairId *big.Int

	Raw *EventMetaData
}

// event Mint(address indexed sender, uint amount0, uint amount1);
type Mint struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int

	Raw *EventMetaData
}

// event Burn(address indexed sender, uint amount0, uint amount1, address indexed to);
type Burn struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	To      common.Address

	Raw *EventMetaData
}

// event Swap(address indexed sender, uint amount0In, uint amount1In, uint amount0Out, uint amount1Out, address indexed to);
type Swap struct {
	Sender     common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To         common.Address

	Raw *EventMetaData
}

// event Sync(uint112 reserve0, uint112 reserve1);
type Sync struct {
	Reserve0 *big.Int
	Reserve1 *big.Int

	Raw *EventMetaData
}
