package speedtestdotnet

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/fopina/speedtest-cli/version"
	"golang.org/x/net/context/ctxhttp"
)

type Client http.Client

type response http.Response

func (client *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err == nil {
		req.Header.Set(
			"User-Agent",
			"Mozilla/5.0 "+
				fmt.Sprintf("(%s; U; %s; en-us)", runtime.GOOS, runtime.GOARCH)+
				fmt.Sprintf("Go/%s", runtime.Version())+
				fmt.Sprintf("(KHTML, like Gecko) speedtest-cli/%s", version.Version))
	}
	return req, err
}

func (c *Client) get(ctx context.Context, url string) (resp *response, err error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	htResp, err := ctxhttp.Do(ctx, (*http.Client)(c), req)
	return (*response)(htResp), err
}

func (c *Client) post(ctx context.Context, url string, bodyType string, body io.Reader) (resp *response, err error) {
	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, body)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = int64(buf.Len())
	htResp, err := ctxhttp.Do(ctx, (*http.Client)(c), req)

	return (*response)(htResp), err
}

func (res *response) ReadContent() ([]byte, error) {
	var content []byte
	if c, err := io.ReadAll(res.Body); err != nil {
		return nil, err
	} else {
		content = c
	}
	if err := res.Body.Close(); err != nil {
		return content, err
	}
	return content, nil
}

func (res *response) ReadXML(out interface{}) error {
	content, err := res.ReadContent()
	if err != nil {
		return err
	}
	return xml.Unmarshal(content, out)
}
