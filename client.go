package ce

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

var (
	baseURL     = mustParseURL("https://api.correct.email")
	emailRegexp = regexp.MustCompile("^[a-z0-9+/_-]+(?:\\.[a-z0-9+/_-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$")
)

// HTTPClient makes all requests to Correct.email API
type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

// Client is handling all requests to Correct.email API
type Client struct {
	apiKey string
	client HTTPClient
}

// New returns new Client
func New(apiKey string) Client {
	return NewWithClient(apiKey, &http.Client{})
}

// NewWithClient returns new Client with custom HTTPClient
func NewWithClient(apiKey string, client HTTPClient) Client {
	return Client{
		apiKey: apiKey,
		client: client,
	}
}

var formDefenderEndpoint = mustParseURL("/v1/single/form-defender/")

// Valid returns whether provided email address and request ip are valid from the
// Form Defender perspective
func (c Client) Valid(email, ip string) (bool, Status, error) {
	// soft regexp before real request
	if !emailRegexp.Match([]byte(email)) {
		return false, Invalid, nil
	}

	query := url.Values{}
	query.Add("key", c.apiKey)
	query.Add("email", email)
	query.Add("ip", ip)

	url := baseURL.ResolveReference(formDefenderEndpoint)
	url.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return false, Invalid, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return false, Invalid, err
	}
	defer response.Body.Close()

	return decodeResponse(response.Body)
}

func decodeResponse(r io.Reader) (bool, Status, error) {
	var resp formDefenderResponse
	err := json.NewDecoder(r).Decode(&resp)
	if err != nil {
		return false, Invalid, err
	}

	if resp.Status == "error" {
		var errorResp errorResponse
		err := json.Unmarshal([]byte(resp.Data), &errorResp)
		if err != nil {
			return false, Invalid, err
		}
		return false, Invalid, applicationErrors[errorResp.Code]
	}

	var validResp validResponse
	return validResp.Valid, validResp.Status, json.Unmarshal(resp.Data, &validResp)
}
