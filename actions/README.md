## 1. How to define an endpoint for standard/registered user?

Consider following example:

```golang
func RegisterRoutes(handlersManager *routes.Manager) {
	group := handlersManager.Group(routes.GroupAPI, "/utxo")
	group.GET("", handlers.AsUser(get))
}
```

-   There is a concept of `RegisterRoutes` func (which is defined for all action's packages, e.g. transactions, utxo, admin, ...)
-   The `RegisterRoutes` function has one required argument `handlersManager *routes.Manager` which is necessary to properly register endpoints and groups.
-   To create a new gin group based on proper prefix and stack of middlewares, use the `handlersManager.Group` method, with desired "group type" and relative path.
-   There are following "group types" available:
    -   `GroupRoot` - with no prefix and only "global" middlewares:
        -   `logging.GinMiddleware(&httpLogger)`
        -   `gin.Recovery()`
        -   `middleware.AppContextMiddleware`
        -   `middleware.CorsMiddleware`
        -   `metrics.requestMetricsMiddleware`
    -   `GroupAPI` with `/api/<api_version> prefix and also `auth_middleware` (besides global middlewares)
    -   `GroupTransactionCallback` with no prefix but with special `middleware.CallbackTokenMiddleware`
-   So... for a new endpoint you'll most probably choose the `GroupAPI`
-   For "registered" user enpoints you should wrap your handler with `handlers.AsUser(yourHandler)`
-   Additionally, the handler func should have following arguments

```golang
func someHandler(c *gin.Context, userContext *reqctx.UserContext)
```

-   Use the `userContext` object to get XPubID of current request.
    -   You can also retrieve the raw XPub by calling `xpub, err := userContext.ShouldGet()`
        -   Keep in mind that this is only possible when authorization is done via `xPub` (not via an `access key`).
        -   For access key authorization, this will return an error.
    -   Additionally, you cannot confuse "user" handlers with others - because of unique set of arguments.

## 2. How to define an endpoint for admin

-   There is no additional middleware for admin (The `auth_middleware` takes care of both `admins` and `"regular" users`)
-   To define a new "admin" endpoint add `handlers.AsAdmin(youAdminHandler)`, as in following example:

```golang
someGroup.POST("/access-keys/count", handlers.AsAdmin(accessKeysCount))
```

-   There is no special "admin" group needed.
-   But admin handlers have unique argument lists:

```golang
func someAdminHandler(c *gin.Context, _ *reqctx.AdminContext)
```

-   `_ *reqctx.AdminContext` doesn't pass any information but helps to distinguish "admin" handlers from other (e.g. "user" handlers or root handlers, like "status")

## 3. How to get "global" objects in a handler

-   The objects are stored in a current request context
-   To get them, use those helper methods:

```golang
func someHandler(c *gin.Context) {
	engine := reqctx.Engine(c)
	logger := reqctx.Logger(c)
	config := reqctx.AppConfig(c)
}
```

## 4. Where the RegisterRoutes functions are called?

-   The `RegisterRoutes` functions (one for each action group) are called in `/actions/register.go` file.
-   If you want to create a new action group you should add the newly created `RegisterRoutes` function into the `Register` function in that file.
