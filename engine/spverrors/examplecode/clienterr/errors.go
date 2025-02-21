package clienterr

import "github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"

var BadRequest = ClientErrorDefinition{
	title:    "Bad request",
	httpCode: 400,
	errType:  errdef.External.NewType("bad_request"),
}

var UnprocessableEntity = ClientErrorDefinition{
	title:    "Unprocessable entity",
	httpCode: 422,
	errType:  errdef.External.NewType("unprocessable_entity"),
}

var Unauthorized = ClientErrorDefinition{
	title:    "Unauthorized",
	httpCode: 401,
	errType:  errdef.External.NewType("unauthorized"),
}
