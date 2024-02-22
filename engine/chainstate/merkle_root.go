package chainstate

import (
	"context"
	"errors"
)

// VerifyMerkleRoots will try to verify merkle roots with all available providers
// When no error is returned, it means that the Block Headers Service client responded with state: Confirmed or UnableToVerify
func (c *Client) VerifyMerkleRoots(ctx context.Context, merkleRoots []MerkleRootConfirmationRequestItem) error {
	pc := c.options.config.blockHedersServiceClient
	if pc == nil {
		c.options.logger.Warn().Msg("VerifyMerkleRoots is called even though no Block Headers Service client is configured; this likely indicates that the paymail capabilities have been cached.")
		return errors.New("no block headers service client found")
	}
	merkleRootsRes, err := pc.verifyMerkleRoots(ctx, c.options.logger, merkleRoots)
	if err != nil {
		return err
	}

	if merkleRootsRes.ConfirmationState == Invalid {
		c.options.logger.Warn().Msg("Not all merkle roots confirmed")
		return errors.New("not all merkle roots confirmed")
	}

	if merkleRootsRes.ConfirmationState == UnableToVerify {
		c.options.logger.Warn().Msg("Some merkle roots were unable to be verified. Proceeding regardless.")
	}

	return nil
}
