# How to Add a new endpoint

## 1. Define the Endpoint in OpenAPI Specification

### Choose the correct API specification file:

- **Admin API** → `api/admin.yaml`
- **User API** → `api/user.yaml`
- **Base API** → `api/base.yaml`

#### Define the endpoint with:

- **Path and method** (e.g., `GET /v2/resource`)
- **Security settings** (`XPubAuth` type and `admin/user` scopes)
- **Request body** (if required)
- **Response body** (or status codes)

### Example of an endpoint definition:

```yaml
  /api/v2/admin/status:
      get:
          security:
              - XPubAuth:
                    - "admin"
          tags:
              - Admin endpoints
          summary: Get admin status
          description: >-
              This endpoint returns admin status. It is used to check if authorization header contain admin xpub.
          responses:
              200:
                  description: Success
              401:
                  $ref: "../components/errors.yaml#/components/responses/NotAuthorized"
```

## 2. Define Security, Request, and Response Bodies

All models should be defined in separate files inside the `/components` directory:

- **Errors** → `/components/errors.yaml`
- **Models** → `/components/models.yaml`
- **Requests** → `/components/requests.yaml`
- **Responses** → `/components/responses.yaml`

> [!NOTE]
> If you want to add endpoint accessible without any authorization just skip the security part in the endpoint definition.

### Example of a request and response definition:

```yaml
components:
  schemas:
    NewRequest:
      type: object
      properties:
        name:
          type: string
        age:
          type: integer
```

## 3. Generate API code

Run the following command to generate the API definition and server interfaces:

```sh
go generate ./...
```

This will generate:

- **API specification** → `gen.api.yaml`
- **Server interface** → `gen.api.go`
- **Models** → `gen.models.go`

## 4. Implement the new endpoint

Server implementation is handled in [/actions/v2/server.go](../actions/v2/server.go) directory

Expand the `Server` struct by:

1. **Adding a new interface** (if needed), like `admin.Server`
2. **Implementing the new endpoint** within an existing interface

### Example

```go
type Server struct {
    admin.Server
    users.Server // new server for user endpoints
}
```

Add your new function implementation to the corresponding interface.

## 5. Run SPV Wallet and test the API

# How to Add a New Middleware

## 1. Create a new middleware

All middlewares should be added in [/server/middleware/](../server/middleware) directory.

Define your middleware function:

```go
func MyNewMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Middleware logic here
        c.Next()
    }
}
```

## 2. Register the middleware in the server handler

Edit [server.go](../server/server.go) to include the new middleware

```go
api.RegisterHandlersWithOptions(ginEngine, v2.NewServer(), api.GinServerOptions{
    BaseURL: "",
    Middlewares: []api.MiddlewareFunc{
        middleware.SignatureAuthWithScopes(),
        // New middleware
    },
    ErrorHandler: func(c *gin.Context, err error, statusCode int) {
        spverrors.ErrorResponse(c, err, log)
    },
})
```

## 3. Run SPV Wallet and test the API

# How to verify SPV Wallet binary do not contain `tools.go`

## What Is `go tool nm`?

The `go tool nm` command **analyzes compiled Go binaries** and lists all symbols (functions, variables, constants) inside them.

- It helps check **which packages and functions are compiled into the final binary**.
- We use it to verify that **code generation tools** (like `oapi-codegen`) are NOT included in the binary.

---

## Running the Check

First, **build the binary**:

```sh
go build -o spvwallet ./cmd/main.go
```

Then, **check if `oapi-codegen` is in the binary**:

```sh
go tool nm spvwallet | grep codegen
```

Expected Output (Good):

```sh
10285ed80 D compress/flate.codegenOrder
101261c28 R github.com/oapi-codegen/runtime..stmp_0
101fbbf40 D github.com/oapi-codegen/runtime/types..inittask
102888638 B github.com/oapi-codegen/runtime/types.emailRegex
100d4fb90 T github.com/oapi-codegen/runtime/types.init
```

### What This Means

- `compress/flate.codegenOrder` → **Not related to `oapi-codegen`** (safe to ignore).
- `github.com/oapi-codegen/runtime/...` → **This is expected! `runtime` is needed for API models.**
- **If you see `github.com/oapi-codegen/oapi-codegen/v2` in the output, then `oapi-codegen` is incorrectly included in the binary.**

---

### When to Run This Check?

Running this test **from time to time** to ensure `oapi-codegen` don’t accidentally get compiled into builds.

- **Before releasing a new version**
- **After updating dependencies (`go mod tidy`)**
- **After modifying build scripts**

To automate this check, you can add/run simple script:

```sh
if go tool nm spvwallet | grep oapi-codegen; then
  echo "❌ ERROR: oapi-codegen is in the binary! Fix required."
  exit 1
else
  echo "✅ SUCCESS: oapi-codegen is NOT in the binary."
fi
```

---

### Final Summary

| Check                                  | Expected Result                                      |
|----------------------------------------|------------------------------------------------------|
| `go tool nm myapp \| grep codegen`      | Only `runtime` references, no `oapi-codegen`        |
| `go list -m all \| grep oapi-codegen`   | `oapi-codegen` should be in `go.mod` (for development only) |
