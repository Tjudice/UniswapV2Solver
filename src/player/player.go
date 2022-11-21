package player

import (
	"UniswapV2Solver/src/data/arango"
	"UniswapV2Solver/src/evts"

	"github.com/ethereum/go-ethereum/common"
)

type Player struct {
	db *arango.DB
}

func NewPlayer(db *arango.DB) *Player {
	return &Player{
		db: db,
	}
}

// Gets section of events
func (p *Player) playEventRange(startBlock, endBlock int, addrs []common.Address) error {
	events := evts.NewEventGroup()
	err := events.AddPairCreatedEvents(p.db, startBlock, endBlock)
	if err != nil {
		return err
	}
	err = events.AddSyncEvents(p.db, startBlock, endBlock)
	if err != nil {
		return err
	}

	return nil
}

// func getEventsInRange(startBlock, endBlock int, addrs []common.Address) ([]evts.UniV2Event, error) {
// 	created
// 	return nil, nil
// }

// const PairCreatedQuery = `
// 	for doc in PairCreatedEvent
// 	filter doc.block >= %d
// 	filter doc.block <= %d
// 	let ret = doc
// `
