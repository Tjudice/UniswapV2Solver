package data

import (
	"fmt"

	"gfx.cafe/open/arango"
	"github.com/arangodb/go-driver"
)

type StageFinished struct {
	Stage    int
	BlockMax int
}

func (a *StageFinished) Key() string {
	return fmt.Sprintf("%d_%d", a.Stage, a.BlockMax)
}

func (a *StageFinished) Collection() (string, []driver.Index) {
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
