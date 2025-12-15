package version

import (
	"fmt"

	"github.com/fopina/speedtest-cli/version"
)

func Main(args []string) {
	fmt.Printf("version: %s\n", version.Version)
}
