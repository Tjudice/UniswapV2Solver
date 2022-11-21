package player

import (
	"UniswapV2Solver/src/data/arango"
	"UniswapV2Solver/src/evts"
	"UniswapV2Solver/src/state"
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Player struct {
	db    *arango.DB
	cl    *ethclient.Client
	state *state.UniswapV2State
}

func NewPlayer(db *arango.DB, cl *ethclient.Client) *Player {
	return &Player{
		db:    db,
		cl:    cl,
		state: state.NewState(cl),
	}
}

func (p *Player) LoadAllPools(ctx context.Context) error {
	blk, err := p.cl.BlockNumber(ctx)
	if err != nil {
		return err
	}
	events := evts.NewEventGroup()
	err = events.AddPairCreatedEvents(ctx, p.db, 0, int(blk))
	if err != nil {
		return err
	}
	events.Sort()
	p.state.AddAllPairs(events.Events)
	return nil
}

// Gets section of events
func (p *Player) PlayEventRange(ctx context.Context, startBlock, endBlock int, addrs []common.Address) error {
	events := evts.NewEventGroup()
	err := events.AddPairCreatedEvents(ctx, p.db, startBlock, endBlock)
	if err != nil {
		return err
	}
	err = events.AddSyncEvents(ctx, p.db, startBlock, endBlock)
	if err != nil {
		return err
	}
	events.Sort()
	for _, evt := range events.Events {
		err = p.state.Update(evt)
		if err != nil {
			return err
		}
	}
	p.state.PrintStateSummary()
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
