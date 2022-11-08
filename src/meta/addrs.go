package meta

import "github.com/ethereum/go-ethereum/common"

const (
	UniswapV2Factory = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
	SushiV2Factory   = "0xC0AEe478e3658e2610c5F7A4A2E1777cE9e4f2Ac"
)

var UniswapV2FactoryAddress = common.HexToAddress(UniswapV2Factory)

var SushiV2FactoryAddress = common.HexToAddress(SushiV2Factory)

var ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
