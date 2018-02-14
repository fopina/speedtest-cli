package speedtest

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	"github.com/jonnrb/speedtest/geo"
	"golang.org/x/sync/errgroup"
)

type ServerID uint64

type Server struct {
	ID          ServerID `xml:"id,attr"`
	Name        string   `xml:"name,attr"`
	Coordinates geo.Coordinates
	URL         string `xml:"url,attr"`
	URL2        string `xml:"url2,attr"`
	Country     string `xml:"country,attr"`
	CC          string `xml:"cc,attr"`
	Sponsor     string `xml:"sponsor,attr"`
	Host        string `xml:"host,attr"`
}

func (s Server) String() string {
	return fmt.Sprintf("%8d: %s (%s, %s) %q", s.ID, s.Sponsor, s.Name, s.Country, s.URL)
}

func (s *Server) RelativeURL(local string) (string, error) {
	u, err := url.Parse(s.URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse server URL %q: %v\n", s.URL, err)
	}

	localURL, err := url.Parse(local)
	if err != nil {
		return "", fmt.Errorf("failed to parse local URL %q: %v\n", local, err)
	}

	return u.ResolveReference(localURL).String(), nil
}

var serverURLs = []string{
	"https://www.speedtest.net/speedtest-servers-static.php",
	"https://c.speedtest.net/speedtest-servers-static.php",
	"https://www.speedtest.net/speedtest-servers.php",
	"https://c.speedtest.net/speedtest-servers.php",
}

func (client *Client) LoadAllServers(ctx context.Context) ([]Server, error) {
	grp, ctx := errgroup.WithContext(ctx)

	c := make(chan []Server)
	for _, url := range serverURLs {
		grp.Go(func() error {
			if s, err := client.loadServersFrom(ctx, url); err != nil {
				return err
			} else {
				c <- s
				return nil
			}
		})
	}

	var servers []Server
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case s := <-c:
				servers = append(servers, s...)
			}
		}
	}()

	if err := grp.Wait(); err != nil {
		return nil, err
	} else {
		return dedupAndSort(servers), nil
	}
}

func (client *Client) loadServersFrom(ctx context.Context, url string) ([]Server, error) {
	resp, err := client.get(ctx, url)
	if resp != nil {
		url = resp.Request.URL.String()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve server list from %q: %v", url, err)
	}

	doc := struct {
		List []struct {
			Server
			Latitude  float64 `xml:"lat,attr"`
			Longitude float64 `xml:"lon,attr"`
		} `xml:"servers>server"`
	}{}

	if err = resp.ReadXML(&doc); err != nil {
		return nil, fmt.Errorf("failed to parse server list from %q: %v", url, err)
	}

	servers := make([]Server, len(doc.List))
	for i, s := range doc.List {
		s.Server.Coordinates = geo.Coordinates{
			Latitude:  geo.Degrees(s.Latitude),
			Longitude: geo.Degrees(s.Longitude),
		}
		servers[i] = s.Server
	}

	return servers, nil
}

func dedupAndSort(servers []Server) []Server {
	m := make(map[ServerID]bool)
	d := make([]Server, len(servers))
	t := 0

	for _, s := range servers {
		if !m[s.ID] {
			d[t] = s
			t += 1
			m[s.ID] = true
		}
	}

	d = d[:t]

	sort.Slice(d, func(i, j int) bool {
		return d[i].ID < d[j].ID
	})

	return d
}
