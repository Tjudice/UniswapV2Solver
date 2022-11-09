package main

import (
	"UniswapV2Solver/src/config"
	"UniswapV2Solver/src/data/arango"
	"UniswapV2Solver/src/grabber"
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
	runner := grabber.NewStageRunner(db, cl, 5000, 10)
	runner.AddStage(grabber.NewStage1(db))
	err = runner.RunStages(context.TODO())
	if err != nil {
		panic(err)
	}
}
