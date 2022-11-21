package data

import (
	"fmt"
	"time"

	"gfx.cafe/open/arango"
	"github.com/arangodb/go-driver"
	"github.com/ethereum/go-ethereum/common"
)

type StageProgress struct {
	Stage int       `json:"stage"`
	Block int       `json:"block"`
	Time  time.Time `json:"time"`
}

func (a *StageProgress) Key() string {
	return fmt.Sprintf("%d_%d", a.Stage, a.Block)
}

func (a *StageProgress) Collection() (string, []driver.Index) {
	return "StageProgress", []driver.Index{
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"block"},
			SetName:   "Block",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"stage"},
			SetName:   "Stage",
			SetUnique: false,
		},
	}
}

type Token struct {
	Address  common.Address `json:"address"`
	Symbol   string         `json:"symbol"`
	Name     string         `json:"name"`
	Decimals uint8          `json:"decimals"`
}

func (a *Token) Key() string {
	return a.Address.Hex()
}

func (a *Token) Collection() (string, []driver.Index) {
	return "Token", []driver.Index{
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"address"},
			SetName:   "Address",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"name"},
			SetName:   "Name",
			SetUnique: false,
		},
		&arango.IDX{
			SetType:   driver.PersistentIndex,
			SetFields: []string{"symbol"},
			SetName:   "Symbol",
			SetUnique: false,
		},
	}
}
