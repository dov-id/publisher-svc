/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type AddFeedbackAttributes struct {
	// course identifier (address) for which feedback is intended
	Course string `json:"course"`
	// feedback itself (its ipfs hash)
	Feedback string `json:"feedback"`
	// network to publish feedback
	Network string `json:"network"`
	// participants’ addresses array
	Participants []string `json:"participants"`
	// list of merkle tree proofs that every address from the addresses list includes is in a corresponding state
	Proofs []SmtProof `json:"proofs"`
	// array that contains public keys of users that will be sign the message and will take part in signing
	PublicKeys []string    `json:"public_keys"`
	Signature  RingSigAttr `json:"signature"`
}
