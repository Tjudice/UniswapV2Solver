package evts

import (
	"UniswapV2Solver/src/data"
	"UniswapV2Solver/src/data/arango"
	"context"
	"fmt"
	"math/big"
	"sort"

	builder "gfx.cafe/open/arango"
)

type EventGroup struct {
	Events []UniV2Event
}

func NewEventGroup() *EventGroup {
	return &EventGroup{
		Events: []UniV2Event{},
	}
}

const PairCreatedEventQuery = `
	for doc in PairCreatedEvent
	filter doc.block >= %d
	filter doc.block <= %d
	let ret = doc
`

func (e *EventGroup) AddPairCreatedEvents(ctx context.Context, db *arango.DB, startBlock, endBlock int) error {
	evts, err := builder.NewBuilder[data.PairCreatedEvent](db.D()).
		Raw(fmt.Sprintf(PairCreatedEventQuery, startBlock, endBlock)).
		Return(ctx, "ret")
	if err != nil {
		return err
	}
	for _, createEvent := range evts {
		mustStr, _ := big.NewInt(0).SetString(createEvent.PairId, 10)
		e.Events = append(e.Events, &PairCreated{
			Token0: createEvent.Token0.Hex(),
			Token1: createEvent.Token1.Hex(),
			Pair:   createEvent.Pair.Hex(),
			PairId: mustStr,
			Raw: &EventMetaData{
				Block:            createEvent.Block,
				Transaction:      createEvent.Transaction,
				TransactionIndex: createEvent.TransactionIndex,
				LogIndex:         createEvent.LogIndex,
				Address:          createEvent.Address.Hex(),
			},
		})
	}
	return nil
}

const SyncEventQuery = `
	for doc in SyncEvent
	filter doc.block >= %d
	filter doc.block <= %d
	let ret = doc
`

func (e *EventGroup) AddSyncEvents(ctx context.Context, db *arango.DB, startBlock, endBlock int) error {
	evts, err := builder.NewBuilder[data.SyncEvent](db.D()).
		Raw(fmt.Sprintf(SyncEventQuery, startBlock, endBlock)).
		Return(ctx, "ret")
	if err != nil {
		return err
	}
	for _, syncEvent := range evts {
		r0, _ := big.NewInt(0).SetString(syncEvent.Reserve0, 10)
		r1, _ := big.NewInt(0).SetString(syncEvent.Reserve1, 10)
		e.Events = append(e.Events, &Sync{
			Reserve0: r0,
			Reserve1: r1,
			Raw: &EventMetaData{
				Block:            syncEvent.Block,
				Transaction:      syncEvent.Transaction,
				TransactionIndex: syncEvent.TransactionIndex,
				LogIndex:         syncEvent.LogIndex,
				Address:          syncEvent.Address,
			},
		})
	}
	return nil
}

func (e *EventGroup) Sort() {
	sort.SliceStable(e.Events, func(i, j int) bool {
		b1, t1, l1 := e.Events[i].EventIndex()
		b2, t2, l2 := e.Events[j].EventIndex()
		if b1 != b2 {
			return b1 < b2
		}
		if t1 != t2 {
			return t1 < t2
		}
		return l1 < l2
	})
}
