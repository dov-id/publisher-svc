/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Feedback struct {
	Key
	Attributes FeedbackAttributes `json:"attributes"`
}
type FeedbackResponse struct {
	Data     Feedback `json:"data"`
	Included Included `json:"included"`
}

type FeedbackListResponse struct {
	Data     []Feedback `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustFeedback - returns Feedback from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustFeedback(key Key) *Feedback {
	var feedback Feedback
	if c.tryFindEntry(key, &feedback) {
		return &feedback
	}
	return nil
}
