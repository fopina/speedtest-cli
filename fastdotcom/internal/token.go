package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

func GetToken(ctx context.Context) (string, error) {
	html, err := getHTML(ctx)
	if err != nil {
		return "", err
	}
	jsPath := extractJSPath(html)
	if jsPath == "" {
		return "", fmt.Errorf("fastdotcom: could not extract fast.com JS URL from the HTML")
	}
	js, err := getJS(ctx, jsPath)
	if err != nil {
		return "", err
	}
	tok := extractToken(js)
	if tok == "" {
		return "", fmt.Errorf("fastdotcom: could not extract fast.com token from JS")
	}
	return tok, nil
}

var c = &http.Client{}

func getHTML(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://fast.com", nil)
	if err != nil {
		return "", fmt.Errorf("fastdotcom: could not get fast.com HTML: %w", err)
	}
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("fastdotcom: could not get fast.com HTML: %w", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("fastdotcom: could not read fast.com HTML: %w", err)
	}
	return string(b), nil
}

var jsRE = regexp.MustCompile("<script.*\"(/app-[[:xdigit:]]+\\.js)\"")

func extractJSPath(html string) string {
	m := jsRE.FindStringSubmatch(html)
	if m == nil {
		return ""
	}
	return m[1]
}

func getJS(ctx context.Context, jsPath string) (string, error) {
	u, err := url.Parse("https://fast.com")
	if err != nil {
		panic(err)
	}
	u.Path = jsPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("fastdotcom: could not get fast.com JS: %w", err)
	}
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("fastdotcom: could not get fast.com JS: %w", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("fastdotcom: could not read fast.com JS: %w", err)
	}
	return string(b), nil
}

var tokenRE = regexp.MustCompile("token:[\"']([[:alpha:]]+)['\"]")

func extractToken(js string) string {
	m := tokenRE.FindStringSubmatch(js)
	if m == nil {
		return ""
	}
	return m[1]
}
