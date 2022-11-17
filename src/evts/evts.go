package evts

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type EventMetaData struct {
	Block            uint64    `json:"block"`
	Transaction      string    `json:"transaction"`
	TransactionIndex uint      `json:"transaction_index"`
	LogIndex         uint      `json:"log_index"`
	Address          string    `json:"address"`
	BlockTimestamp   time.Time `json:"block_timestamp"`
}

// event PairCreated(address indexed token0, address indexed token1, address pair, uint);
type PairCreated struct {
	Token0 string   `json:"token0"`
	Token1 string   `json:"token1"`
	Pair   string   `json:"pair"`
	PairId *big.Int `json:"pair_id"`

	Raw *EventMetaData `json:"raw"`
}

func (e *PairCreated) Name() string {
	return "PairCreated"
}

func (e *PairCreated) EventIndex() (int, int, int) {
	return int(e.Raw.Block), int(e.Raw.TransactionIndex), int(e.Raw.LogIndex)
}

func (e *PairCreated) MetaData() *bind.MetaData {
	return &bind.MetaData{
		ABI: `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"token0","type":"address"},{"indexed":true,"internalType":"address","name":"token1","type":"address"},{"indexed":false,"internalType":"address","name":"pair","type":"address"},{"indexed":false,"internalType":"uint256","name":"","type":"uint256"}],"name":"PairCreated","type":"event"}]`,
	}
}

// event Mint(address indexed sender, uint amount0, uint amount1);
type Mint struct {
	Sender  string   `json:"sender"`
	Amount0 *big.Int `json:"amount0"`
	Amount1 *big.Int `json:"amount1"`

	Raw *EventMetaData `json:"raw"`
}

func (e *Mint) Name() string {
	return "Mint"
}

func (e *Mint) EventIndex() (int, int, int) {
	return int(e.Raw.Block), int(e.Raw.TransactionIndex), int(e.Raw.LogIndex)
}

func (e *Mint) MetaData() *bind.MetaData {
	return &bind.MetaData{
		ABI: `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount0","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1","type":"uint256"}],"name":"Mint","type":"event"}]`,
	}
}

// event Burn(address indexed sender, uint amount0, uint amount1, address indexed to);
type Burn struct {
	Sender  common.Address `json:"sender"`
	Amount0 *big.Int       `json:"amount0"`
	Amount1 *big.Int       `json:"amount1"`
	To      common.Address `json:"to"`

	Raw *EventMetaData `json:"raw"`
}

func (e *Burn) Name() string {
	return "Burn"
}

func (e *Burn) EventIndex() (int, int, int) {
	return int(e.Raw.Block), int(e.Raw.TransactionIndex), int(e.Raw.LogIndex)
}

func (e *Burn) MetaData() *bind.MetaData {
	return &bind.MetaData{
		ABI: `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount0","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1","type":"uint256"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"Burn","type":"event"}]`,
	}
}

// event Swap(address indexed sender, uint amount0In, uint amount1In, uint amount0Out, uint amount1Out, address indexed to);
type Swap struct {
	Sender     common.Address `json:"sender"`
	Amount0In  *big.Int       `json:"amount0in"`
	Amount1In  *big.Int       `json:"amonut1in"`
	Amount0Out *big.Int       `json:"amount0out"`
	Amount1Out *big.Int       `json:"amount1out"`
	To         common.Address `json:"to"`

	Raw *EventMetaData `json:"raw"`
}

func (e *Swap) Name() string {
	return "Swap"
}

func (e *Swap) EventIndex() (int, int, int) {
	return int(e.Raw.Block), int(e.Raw.TransactionIndex), int(e.Raw.LogIndex)
}

func (e *Swap) MetaData() *bind.MetaData {
	return &bind.MetaData{
		ABI: `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount0In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount0Out","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1Out","type":"uint256"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"Swap","type":"event"}]`,
	}
}

// event Sync(uint112 reserve0, uint112 reserve1);
type Sync struct {
	Reserve0 *big.Int `json:"reserve0"`
	Reserve1 *big.Int `json:"reserve1"`

	Raw *EventMetaData `json:"raw"`
}

func (e *Sync) Name() string {
	return "Sync"
}

func (e *Sync) EventIndex() (int, int, int) {
	return int(e.Raw.Block), int(e.Raw.TransactionIndex), int(e.Raw.LogIndex)
}

func (e *Sync) MetaData() *bind.MetaData {
	return &bind.MetaData{
		ABI: `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint112","name":"reserve0","type":"uint112"},{"indexed":false,"internalType":"uint112","name":"reserve1","type":"uint112"}],"name":"Sync","type":"event"}]`,
	}
}
