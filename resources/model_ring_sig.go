/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type RingSig struct {
	Key
	Attributes RingSigAttr `json:"attributes"`
}
type RingSigResponse struct {
	Data     RingSig  `json:"data"`
	Included Included `json:"included"`
}

type RingSigListResponse struct {
	Data     []RingSig `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustRingSig - returns RingSig from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRingSig(key Key) *RingSig {
	var ringSig RingSig
	if c.tryFindEntry(key, &ringSig) {
		return &ringSig
	}
	return nil
}
