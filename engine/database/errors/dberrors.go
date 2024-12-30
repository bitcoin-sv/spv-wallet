package dberrors

import "github.com/bitcoin-sv/spv-wallet/models"

var ErrConvertTOHDPubKey = models.SPVError{Message: "failed to convert to HD public key", StatusCode: 500, Code: "error-convert-to-hd-pubkey"}

var ErrConvertToPubKey = models.SPVError{Message: "failed to convert to public key", StatusCode: 500, Code: "error-convert-to-pubkey"}
