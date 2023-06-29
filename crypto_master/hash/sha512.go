package hash

import "crypto/sha512"

func SHA512StringToString(message string) (digest string) {
	var dig [64]byte = sha512.Sum512([]byte(message))
	digest = string(dig[:])
	return
} //  string -SHA512> string

func SHA512StringToByteArray(message string) (digest [64]byte) {
	digest = sha512.Sum512([]byte(message))
	return
} //  string -SHA512> []byte

func SHA512ByteArrayToString(message []byte) (digest string) {
	var dig [64]byte = sha512.Sum512(message)
	digest = string(dig[:])
	return
} //  []byte -SHA512> string

func SHA512ByteArrayToByteArray(message []byte) (digest [64]byte) {
	digest = sha512.Sum512(message)
	return
} //  []byte -SHA512> []byte
