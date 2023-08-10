/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type AddFeedback struct {
	Key
	Attributes AddFeedbackAttributes `json:"attributes"`
}
type AddFeedbackRequest struct {
	Data     AddFeedback `json:"data"`
	Included Included    `json:"included"`
}

type AddFeedbackListRequest struct {
	Data     []AddFeedback `json:"data"`
	Included Included      `json:"included"`
	Links    *Links        `json:"links"`
}

// MustAddFeedback - returns AddFeedback from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAddFeedback(key Key) *AddFeedback {
	var addFeedback AddFeedback
	if c.tryFindEntry(key, &addFeedback) {
		return &addFeedback
	}
	return nil
}
