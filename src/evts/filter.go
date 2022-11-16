package evts

import (
	"fmt"
	"math/big"
	"strings"

	"gfx.cafe/open/ghost"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ContractGroupFilterer struct {
	abi abi.ABI
	cl  ghost.Client
}

func NewContractGroupFilterer(dat *bind.MetaData, cl ghost.Client) (*ContractGroupFilterer, error) {
	var err error
	o := &ContractGroupFilterer{cl: cl}
	o.abi, err = abi.JSON(strings.NewReader(dat.ABI))
	return o, err
}

func (c *ContractGroupFilterer) FilterLogs(opts *bind.FilterOpts, name string, addrs []common.Address, query ...[]interface{}) (chan types.Log, error) {
	// Don't crash on a lazy user
	if opts == nil {
		opts = &bind.FilterOpts{}
	}
	// Append the event selector to the query parameters and construct the topic set
	query = append([][]interface{}{{c.abi.Events[name].ID}}, query...)

	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return nil, err
	}
	// Start the background filtering
	logs := make(chan types.Log, 128)

	config := ethereum.FilterQuery{
		Addresses: addrs,
		Topics:    topics,
		FromBlock: new(big.Int).SetUint64(opts.Start),
	}
	if opts.End != nil {
		config.ToBlock = new(big.Int).SetUint64(*opts.End)
	}
	buff, err := c.cl.FilterLogs(opts.Context, config)
	if err != nil {
		return nil, err
	}
	for _, log := range buff {
		logs <- log
	}
	return logs, nil
}

func (c *ContractGroupFilterer) UnpackLog(out interface{}, event string, log types.Log) error {
	if log.Topics[0] != c.abi.Events[event].ID {
		return fmt.Errorf("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err := c.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return abi.ParseTopics(out, indexed, log.Topics[1:])
}
