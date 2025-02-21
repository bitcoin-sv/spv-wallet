package errdef

import (
	"errors"
	"fmt"
	"github.com/joomcode/errorx"
)

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func NewProblemDetailsFromError(err error) ProblemDetails {
	var ex *errorx.Error
	if errors.As(err, &ex) {
		title, errType, code := getTitleAndCode(ex)
		if errType == "" {
			errType = ex.Type().FullName()
		}
		return ProblemDetails{
			Type:     errType,
			Title:    title,
			Status:   code,
			Detail:   getDetail(ex),
			Instance: getInstance(ex),
		}
	}

	return ProblemDetails{
		Type:     "unknown_error",
		Title:    "Unknown error",
		Status:   500,
		Instance: "unknown_error",
	}
}

func getInstance(ex *errorx.Error) string {
	instance, ok := ex.Property(PropSpecificProblemOccurrence)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%v", instance)
}

func getDetail(ex *errorx.Error) string {
	var all string
	for _, trait := range globalTraits {
		if ex.HasTrait(trait.Trait) {
			if all != "" {
				all += "; "
			}
			all += trait.Title
		}
	}
	return all
}

func getTitleAndCode(ex *errorx.Error) (title, errType string, code int) {
	clientError, ok := ex.Property(propClientError)
	if ok {
		clientErr := clientError.(ClientError)
		title = clientErr.title
		code = clientErr.httpCode
		errType = clientErr.code
		return
	}

	if ex.IsOfType(UnsupportedOperation) {
		title = "Unsupported Operation"
		code = 501
		return
	}

	title = "Internal Server Error"
	code = 500
	return
}
