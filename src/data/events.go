package data

import (
	"fmt"
	"time"

	"gfx.cafe/open/arango"
	"github.com/arangodb/go-driver"
	"github.com/ethereum/go-ethereum/common"
)

type EventMetaData struct {
	Block            uint64         `json:"block"`
	TransactionIndex uint           `json:"transaction_index"`
	LogIndex         uint           `json:"log_index"`
	Address          common.Address `json:"address"`
	BlockTimestamp   time.Time      `json:"block_timestamp"`
}

// event PairCreated(address indexed token0, address indexed token1, address pair, uint);
type PairCreatedEvent struct {
	Token0 common.Address `json:"token0"`
	Token1 common.Address `json:"token1"`
	Pair   common.Address `json:"pair"`
	PairId string         `json:"pair_id"`

	EventMetaData
}

func (a *PairCreatedEvent) Key() string {
	return fmt.Sprintf("%d_%d_%d", a.Block, a.TransactionIndex, a.LogIndex)
}

func (a *PairCreatedEvent) Collection() (string, []driver.Index) {
	return "PairCreatedEvent", []driver.Index{
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"block"},
			SetName:   "Block",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"pair"},
			SetName:   "Pair",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"pair_id"},
			SetName:   "PairId",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"token0"},
			SetName:   "Token0",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"token1"},
			SetName:   "Token1",
			SetUnique: false,
		},
	}
}

// event Mint(address indexed sender, uint amount0, uint amount1);
type MintEvent struct {
	Sender  common.Address `json:"sender"`
	Amount0 string         `json:"amount0"`
	Amount1 string         `json:"amount1"`

	EventMetaData
}

func (a *MintEvent) Key() string {
	return fmt.Sprintf("%d_%d_%d", a.Block, a.TransactionIndex, a.LogIndex)
}

func (a *MintEvent) Collection() (string, []driver.Index) {
	return "MintEvent", []driver.Index{
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"block"},
			SetName:   "Block",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"sender"},
			SetName:   "Sender",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"address"},
			SetName:   "Pair",
			SetUnique: false,
		},
	}
}

// event Burn(address indexed sender, uint amount0, uint amount1, address indexed to);
type BurnEvent struct {
	Sender  common.Address `json:"sender"`
	Amount0 string         `json:"amount0"`
	Amount1 string         `json:"amount1"`
	To      common.Address `json:"to"`

	EventMetaData
}

func (a *BurnEvent) Key() string {
	return fmt.Sprintf("%d_%d_%d", a.Block, a.TransactionIndex, a.LogIndex)
}

func (a *BurnEvent) Collection() (string, []driver.Index) {
	return "BurnEvent", []driver.Index{
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"block"},
			SetName:   "Block",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"sender"},
			SetName:   "Sender",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"address"},
			SetName:   "Pair",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"to"},
			SetName:   "BurnTo",
			SetUnique: false,
		},
	}
}

// event Swap(address indexed sender, uint amount0In, uint amount1In, uint amount0Out, uint amount1Out, address indexed to);
type SwapEvent struct {
	Sender     common.Address `json:"sender"`
	Amount0In  string         `json:"amount0in"`
	Amount1In  string         `json:"amonut1in"`
	Amount0Out string         `json:"amount0out"`
	Amount1Out string         `json:"amount1out"`
	To         common.Address `json:"to"`

	EventMetaData `json:"raw"`
}

func (a *SwapEvent) Key() string {
	return fmt.Sprintf("%d_%d_%d", a.Block, a.TransactionIndex, a.LogIndex)
}

func (a *SwapEvent) Collection() (string, []driver.Index) {
	return "SwapEvent", []driver.Index{
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"block"},
			SetName:   "Block",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"sender"},
			SetName:   "Sender",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"address"},
			SetName:   "Pair",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"to"},
			SetName:   "SwapTo",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"sender"},
			SetName:   "Sender",
			SetUnique: false,
		},
	}
}

// event Sync(uint112 reserve0, uint112 reserve1);
type SyncEvent struct {
	Reserve0 string `json:"reserve0"`
	Reserve1 string `json:"reserve1"`

	Block            uint64    `json:"block"`
	Transaction      string    `json:"transaction"`
	TransactionIndex uint      `json:"transaction_index"`
	LogIndex         uint      `json:"log_index"`
	Address          string    `json:"address"`
	BlockTimestamp   time.Time `json:"block_timestamp"`
}

func (a *SyncEvent) Key() string {
	return fmt.Sprintf("%d_%d_%d", a.Block, a.TransactionIndex, a.LogIndex)
}

func (a *SyncEvent) Collection() (string, []driver.Index) {
	return "SyncEvent", []driver.Index{
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"block"},
			SetName:   "Block",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"address"},
			SetName:   "Pair",
			SetUnique: false,
		},
	}
}
