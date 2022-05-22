# Definition of business error codes

> The public error code is already in the `github.com/go-eagle/eagle/pkg/errno` package and can be used directly
The error code of the business can be defined by file according to the module

When used, the public error code starts with `errno.`, and the business error code starts with `ecode.`

## Demo

```go
// public error code
import "github.com/go-eagle/eagle/pkg/errno"
...
errno.InternalServerError

// business error code
import "github.com/go-eagle/eagle/internal/ecode"
...
ecode.ErrUserNotFound
```