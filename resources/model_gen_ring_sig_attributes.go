/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type GenRingSigAttributes struct {
	// user index to get its public key from `public_keys` array
	Index int `json:"index"`
	// message that will be signed in ring signature
	Message string `json:"message"`
	// private key of the user for whom `index` parameter was given
	PrivateKey string `json:"private_key"`
	// array that contains public keys of users that will be sign the message and will take part in signing
	PublicKeys []string `json:"public_keys"`
}
