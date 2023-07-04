package helpers

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func StringToBigInt(str string, base int) *big.Int {
	num, _ := new(big.Int).SetString(str, base)
	return num
}

func StringToBytes(str string) []byte {
	if str[:2] == "0x" {
		str = str[2:]
	}
	return common.Hex2Bytes(str)
}

func StringToByte32(str string) [32]byte {
	var result [32]byte
	if str[:2] == "0x" {
		str = str[2:]
	}

	copy(result[:], common.Hex2Bytes(str[:]))
	return result
}

func StringArrToByte32Arr(strs []string) [][32]byte {
	var (
		size   = len(strs)
		result = make([][32]byte, size)
	)

	for i := 0; i < size; i++ {
		if strs[i][:2] == "0x" {
			strs[i] = strs[i][2:]
		}
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
