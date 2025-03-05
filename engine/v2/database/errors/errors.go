package dberrors

import (
	"github.com/joomcode/errorx"
)

var Namespace = errorx.NewNamespace("db")

var QueryFailed = Namespace.NewType("query_failed")
