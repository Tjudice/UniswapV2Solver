package grabber

import "github.com/ethereum/go-ethereum/ethclient"

type Grabber struct {
	c *ethclient.Client
}

func NewGrabber(c *ethclient.Client) (*Grabber, error) {
	return &Grabber{c}, nil
}
