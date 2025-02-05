# How to Add a new endpoint

## 1. Define the Endpoint in OpenAPI Specification

### Choose the correct API specification file:
- **Admin API** → `api/admin.yaml`
- **User API** → `api/user.yaml`
- **Base API** → `api/base.yaml`

#### Define the endpoint with:
- **Path and method** (e.g., `GET /v2/resource`)
- **Security settings** (`XPubAuth` type and `admin/user/basic` scopes)
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
1. **Adding a new interface** (if needed), like `admin.AdminServer`
2. **Implementing the new endpoint** within an existing interface

### Example:
```go
type Server struct {
    admin.AdminServer
    users.UserServer // new server for user endpoints
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
