package grabber

// Filter for all PairCreated events
// Stage2 Must be run after Stage1
// Must log progress so stages can pick up where left off
type Stage1 struct {
	CurrentBlock int
}
