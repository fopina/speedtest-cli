package fastdotcom

import (
	"context"

	"github.com/fopina/speedtest-cli/fastdotcom/internal"
)

type Manifest struct {
	m *internal.Manifest
}

func GetManifest(ctx context.Context, urls int) (*Manifest, error) {
	tok, err := internal.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	mi, err := internal.GetManifest(ctx, tok, urls)
	if err != nil {
		return nil, err
	}
	return &Manifest{mi}, nil
}
