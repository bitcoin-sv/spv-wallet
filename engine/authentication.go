package engine

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// AuthenticateAccessKey check if access key exists
func (c *Client) AuthenticateAccessKey(ctx context.Context, pubAccessKey string) (*AccessKey, error) {
	accessKey, err := getAccessKey(ctx, pubAccessKey, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	} else if accessKey == nil {
		return nil, spverrors.ErrCouldNotFindAccessKey
	} else if accessKey.RevokedAt.Valid {
		return nil, spverrors.ErrAccessKeyRevoked
	}
	return accessKey, nil
}
