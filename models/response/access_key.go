package response

import (
	"time"
)

// AccessKey is a model that represents an access key.
type AccessKey struct {
	// Model is a common model that contains common fields for all models.
	Model
	// ID is an hash of the compressed public key.
	ID string `json:"id" example:"874b86d6fd1d6c85a857e73180164203d8d23211bfd9d04d210f9f7fde5b82d8"`
	// XpubID is an access key's xpub related id.
	XpubID string `json:"xpubId" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50"`
	// RevokedAt is a time when access key was revoked.
	RevokedAt *time.Time `json:"revokedAt,omitempty" example:"2024-02-26T11:02:28.069911Z"`
	// Key is a string representation of an access key.
	Key string `json:"key,omitempty" example:"3fd870d6bf1725f04084cf31209c04be5bd9bed001a390ad3bc632a55a3ee078"`
}
