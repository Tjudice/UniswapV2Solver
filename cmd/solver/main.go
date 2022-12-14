package main

import (
	"UniswapV2Solver/src/config"
	"UniswapV2Solver/src/data/arango"
	"UniswapV2Solver/src/player"
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	c, _, err := config.Load()
	if err != nil {
		panic(err)
	}
	arangoAuth := c.User + ":" + c.Password
	db, err := arango.NewDatabase(c.Host, arangoAuth, "UniswapV2")
	if err != nil {
		panic(err)
	}
	cl, err := ethclient.Dial(c.Rpc)
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()
	player := player.NewPlayer(db, cl)
	err = player.LoadAllPools(ctx)
	if err != nil {
		panic(err)
	}
	err = player.PlayEventRange(context.TODO(), 11000000, 11005000, nil)
	if err != nil {
		panic(err)
	}
}
