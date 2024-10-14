package bhs

// bhsError is an error response that is returned from Block Header Service (BHS)
type bhsError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
