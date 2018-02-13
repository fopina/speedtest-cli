package speedtest

import (
	"context"

	"github.com/jonnrb/speedtest/geo"
)

type Config struct {
	Coordinates        geo.Coordinates
	IP                 string  `xml:"ip,attr"`
	ISP                string  `xml:"isp,attr"`
	ISPRating          float32 `xml:"isprating,attr"`
	ISPDownloadAverage uint    `xml:"ispdlavg,attr"`
	ISPUploadAverage   uint    `xml:"ispulavg,attr"`
	Rating             float32 `xml:"rating,attr"`
}

func (client *Client) Config(ctx context.Context) (Config, error) {
	resp, err := client.get(ctx, "https://www.speedtest.net/speedtest-config.php")
	if err != nil {
		return Config{}, err
	}

	document := struct {
		Client struct {
			Config
			Latitude  float64 `xml:"lat,attr"`
			Longitude float64 `xml:"lon,attr"`
		} `xml:"client"`
	}{}

	if err = resp.ReadXML(&document); err != nil {
		return Config{}, err
	}

	config := document.Client.Config
	config.Coordinates = geo.Coordinates{
		Latitude:  geo.Degrees(document.Client.Latitude),
		Longitude: geo.Degrees(document.Client.Longitude),
	}

	return config, nil
}
