package main

import (
	"UniswapV2Solver/src/grabber"
	"context"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var Config struct {
	RPC_URL string
}

func main() {
	Config.RPC_URL = os.Getenv("RPC_URL")
	cl, err := ethclient.Dial(Config.RPC_URL)
	if err != nil {
		panic(err)
	}
	g, err := grabber.NewGrabber(cl)
	if err != nil {
		panic(err)
	}
	_, err = g.SolveTokenInBlockRange(context.TODO(), common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), 0, 10800000)
	if err != nil {
		panic(err)
	}
}
