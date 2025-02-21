package errors

import "github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"

var ExampleDomainNamespace = errdef.Internal.NewSubNamespace("exampledomain")

var WrongArgument = ExampleDomainNamespace.NewType("wrong_argument", errdef.TraitIllegalArgument)

var SomeARCError = ExampleDomainNamespace.NewType("some_arc_error", errdef.TraitARC)
