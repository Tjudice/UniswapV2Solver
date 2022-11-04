package main

import (
	"UniswapV2Solver/src/grabber"
	"os"

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
}
