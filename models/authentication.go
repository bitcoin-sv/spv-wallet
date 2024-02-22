package models

import "time"

const (
	// AuthHeader is the header to use for authentication (raw xPub)
	AuthHeader = "x-auth-xpub"

	// AuthAccessKey is the header to use for access key authentication (access public key)
	AuthAccessKey = "x-auth-key"

	// AuthSignature is the given signature (body + timestamp)
	AuthSignature = "x-auth-signature"

	// AuthHeaderHash hash of the body coming from the request
	AuthHeaderHash = "x-auth-hash"

	// AuthHeaderNonce random nonce for the request
	AuthHeaderNonce = "x-auth-nonce"

	// AuthHeaderTime the time of the request, only valid for 30 seconds
	AuthHeaderTime = "x-auth-time"

	// AuthSignatureTTL is the max TTL for a signature to be valid
	AuthSignatureTTL = 20 * time.Second
)
