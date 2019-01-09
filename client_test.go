package ce

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type FakeHTTPClient struct {
	resp *http.Response
}

func (c FakeHTTPClient) Do(r *http.Request) (*http.Response, error) {
	return c.resp, nil
}

func TestValid(t *testing.T) {
	testCases := []struct {
		name         string
		resp         *http.Response
		checkedEmail string
		expStatus    Status
		expError     error
		expValid     bool
	}{
		{
			name: "[Success] With clean valid email",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "success",
					"data": {
						"valid":  true,
						"status": "clean"
					}
				}`)),
			},
			checkedEmail: "clean@correct.email",
			expStatus:    Clean,
			expValid:     true,
		},
		{
			name: "[Success] With catch-all valid email",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "success",
					"data": {
						"valid":  true,
						"status": "catch-all"
					}
				}`)),
			},
			checkedEmail: "catchall@correct.email",
			expStatus:    CatchAll,
			expValid:     true,
		},
		{
			name:         "[Success] With syntax incorrect email",
			checkedEmail: ".4123fs@@@fsafcs",
			expStatus:    Invalid,
			expValid:     false,
		},
		{
			name: "[Success] With bounced invalid email",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "success",
					"data": {
						"valid":  false,
						"status": "bounced"
					}
				}`)),
			},
			checkedEmail: "bounce@correct.email",
			expStatus:    Bounced,
			expValid:     false,
		},
		{
			name: "[Success] With special valid email",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "success",
					"data": {
						"valid":  true,
						"status": "special"
					}
				}`)),
			},
			checkedEmail: "special@correct.email",
			expStatus:    Special,
			expValid:     true,
		},
		{
			name: "[Success] With bad MX invalid email",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "success",
					"data": {
						"valid":  false,
						"status": "bad-mx"
					}
				}`)),
			},
			checkedEmail: "badmx@correct.email",
			expStatus:    BadMX,
			expValid:     false,
		},
		{
			name: "[Success] With spam trap invalid email",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "success",
					"data": {
						"valid":  false,
						"status": "spam-trap"
					}
				}`)),
			},
			checkedEmail: "spamtrap@correct.email",
			expStatus:    SpamTrap,
			expValid:     false,
		},
		{
			name: "[Success] With temporary valid email",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "success",
					"data": {
						"valid":  true,
						"status": "temporary"
					}
				}`)),
			},
			checkedEmail: "temporary@correct.email",
			expStatus:    Temporary,
			expValid:     true,
		},
		{
			name: "[Success] With unknown valid email",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "success",
					"data": {
						"valid":  true,
						"status": "unknown"
					}
				}`)),
			},
			checkedEmail: "unknown@correct.email",
			expStatus:    Unknown,
			expValid:     true,
		},

		{
			name: "[Error] With invalid API key",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "error",
					"data": {
						"id":  "app_key_invalid",
						"code": 1000
					}
				}`)),
			},
			checkedEmail: "any@correct.email",
			expError:     applicationErrors[1000],
			expStatus:    Invalid,
		},
		{
			name: "[Error] With blocked account",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "error",
					"data": {
						"id":  "account_blocked",
						"code": 1001
					}
				}`)),
			},
			checkedEmail: "any@correct.email",
			expError:     applicationErrors[1001],
			expStatus:    Invalid,
		},
		{
			name: "[Error] With rate limit exceeded",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "error",
					"data": {
						"id":  "live_rate_limit_reached",
						"code": 1002
					}
				}`)),
			},
			checkedEmail: "any@correct.email",
			expError:     applicationErrors[1002],
			expStatus:    Invalid,
		},
		{
			name: "[Error] With not enough credits",
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{
					"status": "error",
					"data": {
						"id":  "not_enough_credits",
						"code": 1003
					}
				}`)),
			},
			checkedEmail: "any@correct.email",
			expError:     applicationErrors[1003],
			expStatus:    Invalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			client := NewWithClient("", FakeHTTPClient{resp: tt.resp})
			valid, status, err := client.Valid(tt.checkedEmail, "")
			if err != tt.expError {
				t.Errorf("Exp error %v\nGot error %v", tt.expError, err)
				return
			}
			if valid != tt.expValid {
				t.Errorf("Exp valid %v, Got valid %v", tt.expValid, valid)
			}
			if status != tt.expStatus {
				t.Errorf("Exp status %v, Got status %v", tt.expStatus, status)
			}
		})
	}
}
