package UQ112x112

import "math/big"

// library UQ112x112 {
//     uint224 constant Q112 = 2**112;

//     // encode a uint112 as a UQ112x112
//     function encode(uint112 y) internal pure returns (uint224 z) {
//         z = uint224(y) * Q112; // never overflows
//     }

//     // divide a UQ112x112 by a uint112, returning a UQ112x112
//     function uqdiv(uint224 x, uint112 y) internal pure returns (uint224 z) {
//         z = x / uint224(y);
//     }
// }

var q112 = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(112), nil)

func Encode(y *big.Int) *big.Int {
	z := big.NewInt(0).Mul(y, q112)
	return z
}

func Uqdiv(x *big.Int, y *big.Int) *big.Int {
	z := big.NewInt(0).Div(x, y)
	return z
}
