/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type RequestAttributes struct {
	// error in case if request was failed
	Error string `json:"error"`
	// uuid of request to check its status
	Id uuid `json:"id"`
	// the request status, that can be pending, in progress, success and failed
	Status string `json:"status"`
}
