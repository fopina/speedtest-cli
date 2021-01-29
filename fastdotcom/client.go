package fastdotcom

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

func (c *Client) post(ctx context.Context, url string, size int) (*response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, randomBlob(size))
	if err != nil {
		return nil, fmt.Errorf("fastdotcom: could not create request to %q: %w", url, err)
	}
	req.Header.Set("Content-type", "application/octet-stream")
	req.ContentLength = int64(size)
	res, err := (*http.Client)(c).Do(req)
	if err != nil {
		return nil, fmt.Errorf("fastdotcom: could not make request to %q: %w", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"fastdotcom: could not perform upload: status %d %q",
			res.StatusCode, res.Status)
	}
	return (*response)(res), nil
}

func randomBlob(size int) io.Reader {
	return io.LimitReader(rand.Reader, int64(size))
}

func putSizeIntoURL(base string, size int) string {
	u, err := url.Parse(base)
	if err != nil {
		// If we can't parse the URL, it's garbage and we'll fail to actually
		// download the URL.
		return base
	}
	if size < 0 {
		panic("fastdotcom: negative size")
	}
	u.Path += "/range/0-" + strconv.Itoa(size)
	return u.String()
}
