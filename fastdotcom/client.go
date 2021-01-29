package fastdotcom

import (
	"context"
	"fmt"
	"net/http"
)

type Client http.Client

type response http.Response

func (c *Client) get(ctx context.Context, url string) (*response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("fastdotcom: could not create request to %q: %w", url, err)
	}
	res, err := (*http.Client)(c).Do(req)
	if err != nil {
		return nil, fmt.Errorf("fastdotcom: could not make request to %q: %w", url, err)
	}
	return (*response)(res), nil
}
