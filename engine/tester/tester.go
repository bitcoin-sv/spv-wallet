/*
Package tester is a generic testing package with helpful methods for all packages
*/
package tester

import (
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// RandomTablePrefix will make a random prefix (avoid same tables for parallel tests)
func RandomTablePrefix() string {
	prefix, _ := utils.RandomHex(8)
	// add an underscore just in case the table name starts with a number, this is not allowed in sqlite
	return "_" + prefix
}
