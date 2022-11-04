generate:
	mkdir -p generated/uniswap generated/erc20
	./abigen --abi=abi/UniswapV2Pair.abi --type UniswapV2Pair --pkg=uniswap -out=generated/uniswap/uniswapv2pair.go
	./abigen --abi=abi/UniswapV2Factory.abi --type UniswapV2Factory --pkg=uniswap -out=generated/uniswap/uniswapv2factory.go
	./abigen --abi=abi/Erc20.abi --type Erc20Caller --pkg=erc20 -out=generated/erc20/erc20.go




