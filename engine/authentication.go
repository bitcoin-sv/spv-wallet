package engine

import (
	"context"
)

// AuthenticateAccessKey check if access key exists
func (c *Client) AuthenticateAccessKey(ctx context.Context, pubAccessKey string) (*AccessKey, error) {
	accessKey, err := getAccessKey(ctx, pubAccessKey, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	} else if accessKey == nil {
		return nil, ErrUnknownAccessKey
	} else if accessKey.RevokedAt.Valid {
		return nil, ErrAccessKeyRevoked
	}
	return accessKey, nil
}
