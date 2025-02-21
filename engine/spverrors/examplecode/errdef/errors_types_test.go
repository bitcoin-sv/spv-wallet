package errdef

import (
	"fmt"
	"github.com/joomcode/errorx"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestDbConnectionFailed(t *testing.T) {
	err := fmt.Errorf("external error from db package")
	err = DbConnectionFailed.Wrap(err, "wrapped at repo level")
	err = errorx.Decorate(err, "decorated at domain level")

	msg := fmt.Sprintf("%v", err)
	assert.Equal(t, "decorated at domain level, cause: spv-wallet.repo.db_connection_failed: wrapped at repo level, cause: external error from db package", msg)
}
