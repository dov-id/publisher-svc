package hash

import (
	"crypto/sha256"
)

func SHA256StringToString(message string) (digest string) {
	var dig [32]byte = sha256.Sum256([]byte(message))
	digest = string(dig[:])
	return
} // string -SHA256> string

func SHA256StringToByteArray(message string) (digest [32]byte) {
	digest = sha256.Sum256([]byte(message))
	return
} // string -SHA256> []byte

func SHA256ByteArrayToString(message []byte) (digest string) {
	var dig [32]byte = sha256.Sum256(message)
	digest = string(dig[:])
	return
} // []byte -SHA256> string

func SHA256ByteArrayToByteArray(message []byte) (digest [32]byte) {
	digest = sha256.Sum256(message)
	return
} // []byte -SHA256> []byte
