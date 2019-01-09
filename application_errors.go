package ce

import "errors"

var (
	ErrAPIKeyInvalid     = errors.New("Application API key is invalid")
	ErrAccountBlocked    = errors.New("This application owner account is blocked.")
	ErrRateLimitExceeded = errors.New("Rate limit has been exceeded")
	ErrNotEnoughCredits  = errors.New("Account doesn't have enough credits for request")
	ErrServerError       = errors.New("There was an problem with API, please try again")
)

var applicationErrors = map[int]error{
	1000: ErrAPIKeyInvalid,
	1001: ErrAccountBlocked,
	1002: ErrRateLimitExceeded,
	1003: ErrNotEnoughCredits,
	1069: ErrServerError,
}
