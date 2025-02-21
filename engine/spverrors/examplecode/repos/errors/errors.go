package errors

import "github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"

var RepoNamespace = errdef.Internal.NewSubNamespace("repo")

var DbConnectionFailed = RepoNamespace.NewType("db_connection_failed", errdef.TraitConfig)
var DbIllegalArgument = RepoNamespace.NewType("db_illegal_argument", errdef.TraitIllegalArgument)
var DbQueryFailed = RepoNamespace.NewType("db_query_failed")
