package evts

// UniswapV2Event interface
type UniV2Event interface {
	Name() string
	EventIndex() (int, int, int)
	// Filter(addrs []common.Address, start, stop int) error
}
