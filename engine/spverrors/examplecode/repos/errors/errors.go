package errors

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"
	"github.com/joomcode/errorx"
)

var RepoNamespace = errorx.NewNamespace("repo")

var DbConnectionFailed = RepoNamespace.NewType("connection_failed", errdef.TraitConfig)
var DbQueryFailed = RepoNamespace.NewType("query_failed")
var DbShouldNeverHappen = RepoNamespace.NewType("should_never_happen", errdef.TraitShouldNeverHappen)
