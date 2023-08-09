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
	Network string  `json:"network"`
	ZkProof ZkProof `json:"zk_proof"`
}
