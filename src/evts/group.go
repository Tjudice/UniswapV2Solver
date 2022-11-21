package evts

import "UniswapV2Solver/src/data/arango"

type EventGroup struct {
	Events []UniV2Event
}

func NewEventGroup() *EventGroup {
	return &EventGroup{
		Events: []UniV2Event{},
	}
}

func (e *EventGroup) AddPairCreatedEvents(db *arango.DB, startBlock, endBlock int) error {
	return nil
}

func (e *EventGroup) AddSyncEvents(db *arango.DB, startBlock, endBlock int) error {
	return nil
}
