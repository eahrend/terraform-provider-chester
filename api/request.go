package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// doRequest wraps all http requests with the proper authorization
// details.
// On a successful request, it will return a byte slice, and a non-nil error.
// On a failure, it will return a byte slice with details from the API server
// and a non-nil error. Since the REST server should be behind IAP, we'll add the
// proxy-auth header.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	if c.token != nil {
		req.Header.Set("Proxy-Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))
	}
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code: %s", string(b))
	}
	return b, nil
}

// makeRequest is a helper function that builds the http request object
// and then sends it to doRequest.
func (c *Client) makeRequest(b []byte, url, method string) ([]byte, error) {
	var req *http.Request
	var err error
	if b == nil {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(string(b)))
		if err != nil {
			return nil, err
		}
	}
	return c.doRequest(req)
}
