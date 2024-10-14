package chainmodels

// BHSError is an error response that is returned from Block Header Service (BHS)
type BHSError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
