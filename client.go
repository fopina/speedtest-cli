package speedtest

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context/ctxhttp"
)

type Client http.Client

type response http.Response

func (c *Client) get(ctx context.Context, url string) (resp *response, err error) {
	htResp, err := ctxhttp.Get(ctx, (*http.Client)(c), url)
	return (*response)(htResp), err
}

func (c *Client) post(ctx context.Context, url string, bodyType string, body io.Reader) (resp *response, err error) {
	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, &buf)
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
	if c, err := ioutil.ReadAll(res.Body); err != nil {
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
