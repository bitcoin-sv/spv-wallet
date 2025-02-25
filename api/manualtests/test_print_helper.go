package manualtests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type Result interface {
	StatusCode() int
	Bytes() []byte
	Response() *http.Response
}

func Print(result Result) {
	logger := Logger()
	status := result.StatusCode()

	body, fallbackBody := ExtractBody(result)

	response := result.Response()

	url := "<failed to extract>"
	method := "<failed to extract>"
	if response != nil && response.Request != nil {
		url = response.Request.URL.String()
		method = response.Request.Method
	}

	event := logger.Info().
		Str("method", method).
		Str("url", url).
		Int("status", status)

	if body == "" && fallbackBody != nil {
		event.Msgf("============== Response: ==============\n%+v\n", fallbackBody)
	} else {
		event.Msgf("============== Response: ==============\n%s\n", body)
	}
}

func ExtractBody(result Result) (string, any) {
	status := result.StatusCode()
	logger := Logger()

	// if result is pointer, get the value
	if reflect.ValueOf(result).Kind() == reflect.Ptr {
		result = reflect.ValueOf(result).Elem().Interface().(Result)
	}

	val := reflect.ValueOf(result).FieldByName(fmt.Sprintf("JSON%d", status))
	var body string
	var fallbackBody any
	if val.IsValid() {
		fallbackBody = val.Interface()
		jsonBodyB, err := json.MarshalIndent(fallbackBody, "", "  ")
		if err != nil {
			logger.Error().Err(err).Msgf("Cannot marshal json body from field JSON%d", status)
		}
		body = string(jsonBodyB)
	} else {
		body = string(result.Bytes())
	}
	return body, fallbackBody
}

func PrintResponse(t testing.TB, res *http.Response) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Log("Cannot read body:", err)
		bodyBytes = []byte("Cannot read body")
	}

	var jsonBody map[string]any
	if err := json.Unmarshal(bodyBytes, &jsonBody); err != nil {
		t.Log("Cannot unmarshal body:", err)
	}

	response := &resp{
		Body:         bodyBytes,
		Status:       res.StatusCode,
		HTTPResponse: res,
	}

	// put jsonBody into response field of name JSON + res.StatusCode
	if jsonBody != nil {
		responseField := reflect.ValueOf(response).Elem().FieldByName(fmt.Sprintf("JSON%d", res.StatusCode))
		if responseField.IsValid() {
			responseField.Set(reflect.ValueOf(jsonBody))
		}
	}

	Print(response)
}

type resp struct {
	Body         []byte
	Status       int
	JSON200      map[string]any
	JSON201      map[string]any
	JSON400      map[string]any
	JSON401      map[string]any
	JSON404      map[string]any
	JSON422      map[string]any
	JSON500      map[string]any
	HTTPResponse *http.Response
}

func (r resp) StatusCode() int {
	return r.Status
}

func (r resp) Response() *http.Response {
	return r.HTTPResponse
}

func (r resp) Bytes() []byte {
	return r.Body
}
