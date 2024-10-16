package apierror

import "fmt"

const MissingAuthHeaderJSON = `{ "code":"error-unauthorized-auth-header-missing", "message":"missing auth header" }`
const AdminNotAuthorizedJSON = `{ "code":"error-admin-auth-on-user-endpoint", "message":"cannot call user's endpoints with admin authorization" }`
const CannotBindBodyJSON = `{"code":"error-bind-body-invalid", "message":"cannot bind request body"}`

func ExpectedJSON(code string, message string) string {
	return fmt.Sprintf(`{"code":"%s","message":"%s"}`, code, message)
}