//nolint:revive // Error types should be self-explanatory
package clienterr

var BadRequest = ClientErrorDefinition{
	title:    "Bad request",
	typeName: "bad_request",
	httpCode: 400,
}

var UnprocessableEntity = ClientErrorDefinition{
	title:    "Unprocessable entity",
	typeName: "unprocessable_entity",
	httpCode: 422,
}
