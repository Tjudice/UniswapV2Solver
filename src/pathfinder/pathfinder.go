package path

import "UniswapV2Solver/src/state"

type PathFinder struct {
	state *state.UniswapV2State
}

func NewPathFinder(state *state.UniswapV2State) *PathFinder {
	return &PathFinder{
		state: state,
	}
}
