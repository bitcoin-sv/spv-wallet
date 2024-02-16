package models

import "time"

const (
	// AuthHeader is the header to use for authentication (raw xPub)
	AuthHeader = "spv-wallet-auth-xpub"

	// AuthAccessKey is the header to use for access key authentication (access public key)
	AuthAccessKey = "spv-wallet-auth-key"

	// AuthSignature is the given signature (body + timestamp)
	AuthSignature = "spv-wallet-auth-signature"

	// AuthHeaderHash hash of the body coming from the request
	AuthHeaderHash = "spv-wallet-auth-hash"

	// AuthHeaderNonce random nonce for the request
	AuthHeaderNonce = "spv-wallet-auth-nonce"

	// AuthHeaderTime the time of the request, only valid for 30 seconds
	AuthHeaderTime = "spv-wallet-auth-time"

	// AuthSignatureTTL is the max TTL for a signature to be valid
	AuthSignatureTTL = 20 * time.Second
)
