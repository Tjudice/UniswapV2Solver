package evts

import "github.com/ethereum/go-ethereum/accounts/abi/bind"

// UniswapV2Event interface
type UniV2Event interface {
	Name() string
	EventIndex() (int, int, int)
	MetaData() *bind.MetaData
	// Filter(addrs []common.Address, start, stop int) error
}
