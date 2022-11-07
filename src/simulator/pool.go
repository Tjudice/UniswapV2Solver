package simulator

import (
	"UniswapV2Solver/src/evts"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Position struct {
}

type Pool struct {
	PairId int `json:"pair_id"`

	Token0 common.Address `json:"token0"`
	Token1 common.Address `json:"token1"`

	PoolAddress common.Address `json:"pool_address"`

	Reserve0 *big.Int `json:"reserve0"`
	Reserve1 *big.Int `json:"reserve1"`

	Balance0 *big.Int `json:"balance0"`
	Balance1 *big.Int `json:"balance1"`

	Price0CumulativeLast *big.Int `json:"price0_cumulative_last"`
	Price1CumulativeLast *big.Int `json:"price1_cumulative_last"`
	KLast                *big.Int `json:"k_last"`

	BlockTimestampLast time.Time `json:"block_timestamp_last"`

	CreateEvent *evts.PairCreated `json:"create_event"`

	Positions map[common.Address]Position `json:"positions"`
}

func NewPool(evt *evts.PairCreated) *Pool {
	return &Pool{
		PairId:               int(evt.PairId.Int64()),
		Token0:               evt.Token0,
		Token1:               evt.Token1,
		PoolAddress:          evt.Pair,
		Reserve0:             big.NewInt(0),
		Reserve1:             big.NewInt(0),
		Balance0:             big.NewInt(0),
		Balance1:             big.NewInt(0),
		Price0CumulativeLast: big.NewInt(0),
		Price1CumulativeLast: big.NewInt(0),
		KLast:                big.NewInt(0),
		BlockTimestampLast:   evt.Raw.BlockTimestamp,
		CreateEvent:          evt,
	}
}

func (p *Pool) SimulateEvent(evt evts.UniV2Event) error {
	switch e := evt.(type) {
	case *evts.Mint:
		return p.Mint(e)
	case *evts.Burn:
		return p.Burn(e)
	case *evts.Swap:
		return p.Swap(e)
	case *evts.Sync:
		return p.Sync(e)
	default:
		return fmt.Errorf("unknown event type: %T", evt)
	}
}

// function mint(address to) external lock returns (uint liquidity) {
// 	(uint112 _reserve0, uint112 _reserve1,) = getReserves(); // gas savings
// 	uint balance0 = IERC20(token0).balanceOf(address(this));
// 	uint balance1 = IERC20(token1).balanceOf(address(this));
// 	uint amount0 = balance0.sub(_reserve0);
// 	uint amount1 = balance1.sub(_reserve1);

// 	bool feeOn = _mintFee(_reserve0, _reserve1);
// 	uint _totalSupply = totalSupply; // gas savings, must be defined here since totalSupply can update in _mintFee
// 	if (_totalSupply == 0) {
// 		liquidity = Math.sqrt(amount0.mul(amount1)).sub(MINIMUM_LIQUIDITY);
// 	   _mint(address(0), MINIMUM_LIQUIDITY); // permanently lock the first MINIMUM_LIQUIDITY tokens
// 	} else {
// 		liquidity = Math.min(amount0.mul(_totalSupply) / _reserve0, amount1.mul(_totalSupply) / _reserve1);
// 	}
// 	require(liquidity > 0, 'UniswapV2: INSUFFICIENT_LIQUIDITY_MINTED');
// 	_mint(to, liquidity);

// 	_update(balance0, balance1, _reserve0, _reserve1);
// 	if (feeOn) kLast = uint(reserve0).mul(reserve1); // reserve0 and reserve1 are up-to-date
// 	emit Mint(msg.sender, amount0, amount1);
// }

func (p *Pool) Mint(e *evts.Mint) error {
	return nil
}

// function burn(address to) external lock returns (uint amount0, uint amount1) {
// 	(uint112 _reserve0, uint112 _reserve1,) = getReserves(); // gas savings
// 	address _token0 = token0;                                // gas savings
// 	address _token1 = token1;                                // gas savings
// 	uint balance0 = IERC20(_token0).balanceOf(address(this));
// 	uint balance1 = IERC20(_token1).balanceOf(address(this));
// 	uint liquidity = balanceOf[address(this)];

// 	bool feeOn = _mintFee(_reserve0, _reserve1);
// 	uint _totalSupply = totalSupply; // gas savings, must be defined here since totalSupply can update in _mintFee
// 	amount0 = liquidity.mul(balance0) / _totalSupply; // using balances ensures pro-rata distribution
// 	amount1 = liquidity.mul(balance1) / _totalSupply; // using balances ensures pro-rata distribution
// 	require(amount0 > 0 && amount1 > 0, 'UniswapV2: INSUFFICIENT_LIQUIDITY_BURNED');
// 	_burn(address(this), liquidity);
// 	_safeTransfer(_token0, to, amount0);
// 	_safeTransfer(_token1, to, amount1);
// 	balance0 = IERC20(_token0).balanceOf(address(this));
// 	balance1 = IERC20(_token1).balanceOf(address(this));

// 	_update(balance0, balance1, _reserve0, _reserve1);
// 	if (feeOn) kLast = uint(reserve0).mul(reserve1); // reserve0 and reserve1 are up-to-date
// 	emit Burn(msg.sender, amount0, amount1, to);
// }

func (p *Pool) Burn(e *evts.Burn) error {
	return nil
}

// function swap(uint amount0Out, uint amount1Out, address to, bytes calldata data) external lock {
// 	require(amount0Out > 0 || amount1Out > 0, 'UniswapV2: INSUFFICIENT_OUTPUT_AMOUNT');
// 	(uint112 _reserve0, uint112 _reserve1,) = getReserves(); // gas savings
// 	require(amount0Out < _reserve0 && amount1Out < _reserve1, 'UniswapV2: INSUFFICIENT_LIQUIDITY');

// 	uint balance0;
// 	uint balance1;
// 	{ // scope for _token{0,1}, avoids stack too deep errors
// 	address _token0 = token0;
// 	address _token1 = token1;
// 	require(to != _token0 && to != _token1, 'UniswapV2: INVALID_TO');
// 	if (amount0Out > 0) _safeTransfer(_token0, to, amount0Out); // optimistically transfer tokens
// 	if (amount1Out > 0) _safeTransfer(_token1, to, amount1Out); // optimistically transfer tokens
// 	if (data.length > 0) IUniswapV2Callee(to).uniswapV2Call(msg.sender, amount0Out, amount1Out, data);
// 	balance0 = IERC20(_token0).balanceOf(address(this));
// 	balance1 = IERC20(_token1).balanceOf(address(this));
// 	}
// 	uint amount0In = balance0 > _reserve0 - amount0Out ? balance0 - (_reserve0 - amount0Out) : 0;
// 	uint amount1In = balance1 > _reserve1 - amount1Out ? balance1 - (_reserve1 - amount1Out) : 0;
// 	require(amount0In > 0 || amount1In > 0, 'UniswapV2: INSUFFICIENT_INPUT_AMOUNT');
// 	{ // scope for reserve{0,1}Adjusted, avoids stack too deep errors
// 	uint balance0Adjusted = balance0.mul(1000).sub(amount0In.mul(3));
// 	uint balance1Adjusted = balance1.mul(1000).sub(amount1In.mul(3));
// 	require(balance0Adjusted.mul(balance1Adjusted) >= uint(_reserve0).mul(_reserve1).mul(1000**2), 'UniswapV2: K');
// 	}

// 	_update(balance0, balance1, _reserve0, _reserve1);
// 	emit Swap(msg.sender, amount0In, amount1In, amount0Out, amount1Out, to);
// }

func (p *Pool) Swap(e *evts.Swap) error {
	return nil
}

// function _update(uint balance0, uint balance1, uint112 _reserve0, uint112 _reserve1) private {
// 	require(balance0 <= uint112(-1) && balance1 <= uint112(-1), 'UniswapV2: OVERFLOW');
// 	uint32 blockTimestamp = uint32(block.timestamp % 2**32);
// 	uint32 timeElapsed = blockTimestamp - blockTimestampLast; // overflow is desired
// 	if (timeElapsed > 0 && _reserve0 != 0 && _reserve1 != 0) {
// 		// * never overflows, and + overflow is desired
// 		price0CumulativeLast += uint(UQ112x112.encode(_reserve1).uqdiv(_reserve0)) * timeElapsed;
// 		price1CumulativeLast += uint(UQ112x112.encode(_reserve0).uqdiv(_reserve1)) * timeElapsed;
// 	}
// 	reserve0 = uint112(balance0);
// 	reserve1 = uint112(balance1);
// 	blockTimestampLast = blockTimestamp;
// 	emit Sync(reserve0, reserve1);
// }

func (p *Pool) Sync(e *evts.Sync) error {
	p.Reserve0.Set(e.Reserve0)
	p.Reserve1.Set(e.Reserve1)
	p.BlockTimestampLast = e.Raw.BlockTimestamp
	return nil
}
