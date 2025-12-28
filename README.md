# contest

`contest` is a lightweight DSL for testing [ConnectRPC](https://connectrpc.com) handlers.
It executes handlers directly using Go's `httptest` package, automatically decodes `connect.Error` responses, and
provides a fluent, expressive API for asserting response behavior.

---

## Features

* Runs Connect handlers directly without a network layer
* Automatically unmarshals `connect.Error` JSON responses
* Provides fluent assertion methods:
    * `ExpectStatus`
    * `ExpectHeader`
    * `Out` (for unmarshalling protobuf responses)
    * `Err` (returns `*connect.Error` if present)
* Supports structured error detail decoding (`connect.ErrorDetail`)
* Compatible with `connect.WithInterceptors(...)`

---

## Example

```go
package handler_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/contest"
	authv1 "github.com/mickamy/contest/example/gen/auth/v1"
	"github.com/mickamy/contest/example/gen/github.com/mickamy/contest/example/gen/auth/v1/authv1connect"
	"github.com/mickamy/contest/example/internal/domain/session/handler"
)

func TestSession_SignIn(t *testing.T) {
	t.Parallel()

	var out authv1.SignInResponse

	_, server := authv1connect.NewSessionServiceHandler(handler.NewSession())
	ct := contest.
		New(t, server).
		Procedure(authv1connect.SessionServiceSignInProcedure).
		In(&authv1.SignInRequest{
			Email:    "user@example.com",
			Password: "password123",
		}).
		Do().
		ExpectStatus(http.StatusOK)

	ct.Out(&out)
	require.NotZero(t, out.Tokens)
	assert.NotZero(t, out.Tokens.Access)
	assert.NotZero(t, out.Tokens.Refresh)
}
```

---

## API Overview

### Creating a client

```go
ct := contest.New(t, handler)
```

### Configuring a request

```go
ct.Procedure("/package.Service/Method").Header("X-Test", "1").In(&req)
```

### Executing and asserting

```go
ct.Do().ExpectStatus(http.StatusOK).Out(&res)
```

### Handling errors

```go
if connErr := ct.Err(); connErr != nil {
    fmt.Println(connErr.Code(), connErr.Message())
}
```

---

## License

[MIT](./LICENSE)
