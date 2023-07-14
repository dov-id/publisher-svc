package helpers

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func Trim0xPrefix(str string) string {
	if str[:2] == "0x" {
		str = str[2:]
	}

	return str
}

func StringToBigInt(str string, base int) *big.Int {
	num, _ := new(big.Int).SetString(str, base)
	return num
}

func StringToBytes(str string) []byte {
	str = Trim0xPrefix(str)

	return common.Hex2Bytes(str)
}

func StringToByte32(str string) [32]byte {
	var result [32]byte
	str = Trim0xPrefix(str)

	copy(result[:], common.Hex2Bytes(str[:]))
	return result
}

func StringArrToByte32Arr(strs []string) [][32]byte {
	var (
		size   = len(strs)
		result = make([][32]byte, size)
	)

	for i := 0; i < size; i++ {
		strs[i] = Trim0xPrefix(strs[i])
		result[i] = StringToByte32(strs[i])
	}

	return result
}

func StringArrToBigIntArr(strs []string, base int) []*big.Int {
	var (
		size   = len(strs)
		result = make([]*big.Int, size)
	)

	for i := 0; i < size; i++ {
		result[i] = StringToBigInt(strs[i], base)
	}

	return result
}
