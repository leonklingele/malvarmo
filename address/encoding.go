package address

import (
	"math/big"
)

// Based on https://github.com/moneromooo-monero/monero-wallet-generator/blob/master/monero-wallet-generator.html
// base58encode converts data into Base58-format
func base58encode(data []byte) []byte {
	const (
		fullBlockSize        = 8
		fullEncodedBlockSize = 11
	)
	var (
		alphabet          = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
		encodedBlockSizes = []int{0, 2, 3, 5, 6, 7, 9, 10, 11}
	)

	encodeBlock := func(data, buf []byte, index int) []byte {
		lenAlphabet := big.NewInt(int64(len(alphabet)))
		num := big.NewInt(0).SetBytes(data)
		rem := big.NewInt(0)

		i := encodedBlockSizes[len(data)] - 1
		for num.Cmp(big.NewInt(0)) == 1 {
			num.QuoRem(num, lenAlphabet, rem)
			buf[index+i] = alphabet[int(rem.Int64())]
			i--
		}

		return buf
	}

	fullBlockCount := len(data) / fullBlockSize
	lastBlockSize := len(data) % fullBlockSize
	resSize := fullBlockCount*fullEncodedBlockSize + encodedBlockSizes[lastBlockSize]
	res := make([]byte, resSize)
	for i := 0; i < resSize; i++ {
		res[i] = alphabet[0]
	}

	for i := 0; i < fullBlockCount; i++ {
		res = encodeBlock(
			data[i*fullBlockSize:i*fullBlockSize+fullBlockSize],
			res,
			i*fullEncodedBlockSize,
		)
	}

	if lastBlockSize > 0 {
		res = encodeBlock(
			data[fullBlockCount*fullBlockSize:fullBlockCount*fullBlockSize+lastBlockSize],
			res,
			fullBlockCount*fullEncodedBlockSize,
		)
	}

	return res
}
