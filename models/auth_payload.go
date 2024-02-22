package models

// AuthPayload is the struct that is used to create the signature for the API call
type AuthPayload struct {
	// AuthHash is the hash of the body contents
	AuthHash string `json:"auth_hash"`
	// AuthNonce is a random string
	AuthNonce string `json:"auth_nonce"`
	// AuthTime is the current time in milliseconds
	AuthTime int64 `json:"auth_time"`
	// BodyContents is the body of the request
	BodyContents string `json:"body_contents"`
	// Signature is the signature of the body contents
	Signature string `json:"signature"`
	// XPub is the xpub of the account
	XPub string `json:"xpub"`
	// AccessKey is the access key of the account
	AccessKey string `json:"access_key"`
}
