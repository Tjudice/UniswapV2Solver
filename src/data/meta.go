package data

import (
	"fmt"
	"time"

	"gfx.cafe/open/arango"
	"github.com/arangodb/go-driver"
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
