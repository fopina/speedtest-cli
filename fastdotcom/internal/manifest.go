package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Manifest struct {
	Client  ManifestClient   `json:"client"`
	Targets []ManifestTarget `json:"targets"`
}

type ManifestClient struct {
	ASN      string           `json:"asn"`
	ISP      string           `json:"isp"`
	Location ManifestLocation `json:"location"`
	IP       string           `json:"ip"`
}

type ManifestTarget struct {
	Name     string           `json:"name"`
	URL      string           `json:"url"`
	Location ManifestLocation `json:"location"`
}

type ManifestLocation struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func GetManifest(ctx context.Context, token string, urls int) (*Manifest, error) {
	ms, err := getManifestString(ctx, makeManifestURL(token, urls))
	if err != nil {
		return nil, err
	}
	return unmarshalManifest(ms)
}

func unmarshalManifest(s string) (*Manifest, error) {
	var m Manifest
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, fmt.Errorf("fastdotcom: couldn't unmarshal manifest: %w", err)
	}
	return &m, nil
}

func getManifestString(ctx context.Context, u string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("fastdotcom: could not create manifest request: %w", err)
	}
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("fastdotcom: could not get manifest: %w", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("fastdotcom: could not read manifest: %w", err)
	}
	return string(b), nil
}

func makeManifestURL(token string, urls int) string {
	u, err := url.Parse("https://api.fast.com/netflix/speedtest/v2")
	if err != nil {
		panic(err)
	}

	q := make(url.Values)
	q.Set("https", "true")
	q.Set("token", token)
	q.Set("urlCount", strconv.Itoa(urls))

	u.RawQuery = q.Encode()
	return u.String()
}
