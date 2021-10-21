// Package API is a client side SDK to the chester-api http server
package api

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"net/http"
)

// ClientOption is an option wrapper for the client
type ClientOption func(*Client)

// Client is the wrapper for all the things the
// client API will need.
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	token      *oauth2.Token
	Username   string
	Password   string
	audience   string
}

// NewClient creates a pointer to a Client struct with specific
// configuration options.
// On a failure it will return a nil object and a non-nil error.
func NewClient(host, user, pass, audience string) (*Client, error) {
	ctx := context.Background()
	if audience == "" {
		return nil, fmt.Errorf("no audience found")
	}
	ts, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("failed to create token from audience %s, error: %s", audience, err.Error())
	}
	c := &Client{
		HTTPClient: &http.Client{},
	}
	if host == "" {
		return nil, fmt.Errorf("no host found")
	}
	if user == "" {
		return nil, fmt.Errorf("no user found")
	}
	if pass == "" {
		return nil, fmt.Errorf("no pass found")
	}
	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	c.token = token
	c.HostURL = host
	c.Username = user
	c.Password = pass

	return c, nil
}

// NewClientWithOptions allows for users to create a client with their own options.
func NewClientWithOptions(opts ...ClientOption) (*Client, error) {
	c := &Client{
		HTTPClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.audience != "" && c.token == nil {
		ts, err := idtoken.NewTokenSource(context.Background(), c.audience)
		if err != nil {
			return nil, err
		}
		token, err := ts.Token()
		if err != nil {
			return nil, err
		}
		c.token = token
	}
	return c, nil
}

// WithHost creates a ClientOption that modifies the client's
// base host url
func WithHost(host string) ClientOption {
	return func(c *Client) {
		c.HostURL = host
	}
}

// WithUsername creates a ClientOption that modifies the
// client's basic auth username
func WithUsername(username string) ClientOption {
	return func(c *Client) {
		c.Username = username
	}
}

// WithPassword creates a ClientOption that modifies
// the client's basic auth password
func WithPassword(password string) ClientOption {
	return func(c *Client) {
		c.Password = password
	}
}

// WithToken creates a ClientOption that sets the
// the client's oauth2 token.
//
// !!! DO NOT USE IN CONJUNCTION WITH WithAudience!!!
func WithToken(token *oauth2.Token) ClientOption {
	return func(c *Client) {
		c.token = token
	}
}

// WithAudience sets the client's audience and generates
// a token using google's idtoken package.
//
// !!! DO NOT USE IN CONJUNCTION WITH WithToken !!!
func WithAudience(audience string) ClientOption {
	return func(c *Client) {
		c.audience = audience
	}
}

// WithHTTPClient creates a ClientOption that overrides the
// default http client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = client
	}
}
