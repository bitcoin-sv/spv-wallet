//go:build database_tests
// +build database_tests

package engine

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestTestSuite kick-starts all suite tests
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(EmbeddedDBTestSuite))
}
