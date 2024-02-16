package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestAccessKey tests AccessKey model.
func TestAccessKey(t *testing.T) {
	ac := new(AccessKey)
	ac.Model.UpdatedAt = time.Now().UTC()
	ac.Model.CreatedAt = time.Now().UTC()
	ac.Model.DeletedAt = time.Now().UTC()
	ac.XpubID = "123"
	ac.ID = "123"

	require.Equal(t, "123", ac.ID)
}

// ExampleAccessKey is an example for AccessKey model.
func ExampleAccessKey() {
	ac := new(AccessKey)
	ac.Model.UpdatedAt = time.Now().UTC()
	ac.Model.CreatedAt = time.Now().UTC()
	ac.Model.DeletedAt = time.Now().UTC()
	ac.XpubID = "123"
	ac.ID = "123"
	fmt.Printf("%s", ac.ID)
	// Output: 123
}
