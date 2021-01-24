package units

import (
	"fmt"
)

type BytesPerSecond float64
type BitsPerSecond float64

const (
	KBps BytesPerSecond = 1000
	MBps                = 1000 * KBps
	GBps                = 1000 * MBps

	Kbps BitsPerSecond = 1000
	Mbps               = 1000 * Kbps
	Gbps               = 1000 * Mbps
)

func (s BytesPerSecond) BitsPerSecond() BitsPerSecond {
	return BitsPerSecond(float64(s) * 8)
}

func (s BytesPerSecond) String() string {
	if s < KBps {
		return fmt.Sprintf("%.0f B/s", s)
	} else if s < MBps {
		return fmt.Sprintf("%.02f KB/s", s/KBps)
	} else if s < GBps {
		return fmt.Sprintf("%.02f MB/s", s/MBps)
	} else {
		return fmt.Sprintf("%.02f GB/s", s/GBps)
	}
}

func (s BitsPerSecond) BytesPerSecond() BytesPerSecond {
	return BytesPerSecond(float64(s) / 8)
}

func (s BitsPerSecond) String() string {
	if s < Kbps {
		return fmt.Sprintf("%.0f b/s", s)
	} else if s < Mbps {
		return fmt.Sprintf("%.02f Kb/s", s/Kbps)
	} else if s < Gbps {
		return fmt.Sprintf("%.02f Mb/s", s/Mbps)
	} else {
		return fmt.Sprintf("%.02f Gb/s", s/Gbps)
	}
}
