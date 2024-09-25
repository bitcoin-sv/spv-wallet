package merkleroots

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type BHSErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// mapBHSErrorResponseToSpverror is a method that will check what kind of response came back from
// Block Header Service and map it to spverror and set it to context
func mapBHSErrorResponseToSpverror(ctx *gin.Context, res *http.Response, logger *zerolog.Logger) {
	var responseErr BHSErrorResponse

	err := json.NewDecoder(res.Body).Decode(&responseErr)
	if err != nil {
		spverrors.ErrorResponse(ctx, spverrors.ErrBHSParsingResponse, logger)
		return
	}

	err = models.SPVError{
		Message:    responseErr.Message,
		StatusCode: res.StatusCode,
		Code:       mapBHSCodeToSpverrorCode(responseErr.Code),
	}
	spverrors.ErrorResponse(ctx, err, logger)
}

// mapBHSCodeToSpverrorCode maps error code returned from Block Header Service to
// match error codes structure defined in spverrors
func mapBHSCodeToSpverrorCode(code string) string {
	if code == "" {
		return models.UnknownErrorCode
	}

	// check if code starts with "Err" followed by an uppercase letter
	if len(code) > 3 && strings.HasPrefix(code, "Err") && unicode.IsUpper(rune(code[3])) {
		code = strings.Replace(code, "Err", "error", 1)
	}

	var result strings.Builder

	for i, char := range code {
		// if it's the first character and uppercase, make it lowercase
		if i == 0 && unicode.IsUpper(char) {
			result.WriteRune(unicode.ToLower(char))
		} else if unicode.IsUpper(char) {
			// If it's uppercase, prepend with a hypen and make it lowercase
			result.WriteRune('-')
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}

	}

	return result.String()
}
