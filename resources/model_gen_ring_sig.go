/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type GenRingSig struct {
	Key
	Attributes GenRingSigAttributes `json:"attributes"`
}
type GenRingSigResponse struct {
	Data     GenRingSig `json:"data"`
	Included Included   `json:"included"`
}

type GenRingSigListResponse struct {
	Data     []GenRingSig `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustGenRingSig - returns GenRingSig from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustGenRingSig(key Key) *GenRingSig {
	var genRingSig GenRingSig
	if c.tryFindEntry(key, &genRingSig) {
		return &genRingSig
	}
	return nil
}
