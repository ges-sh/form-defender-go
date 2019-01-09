### Go client for Correct.email Form Defender

#### What is Form Defender?
https://correct.email/form-defender/

### How to install:
`go get github.com/ges-sh/form-defender-go`

### How to use:
```go
package main

import (
	"fmt"

	ce "github.com/ges-sh/form-defender-go"
)

func main() {
	ceClient := ce.New("apiKey")
	valid, status, err := ceClient.Valid("test@correct.email", "0.0.0.0")
	if err != nil {
		switch err {
		case ce.ErrAPIKeyInvalid:
		case ce.ErrAccountBlocked:
		case ce.ErrRateLimitExceeded:
		case ce.ErrNotEnoughCredits:
		default:
		}
	}

	fmt.Println(valid, status)
}
```
