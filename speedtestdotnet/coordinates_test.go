package speedtestdotnet

import (
	"math/rand"
	"testing"
	"time"

	"go.jonnrb.io/speedtest/geo"
)

func TestSortServersByDistance(t *testing.T) {
	expected := []Server{
		Server{
			ID: ServerID(0),
			Coordinates: geo.Coordinates{
				Latitude:  geo.Degrees(0),
				Longitude: geo.Degrees(0),
			},
		},
		Server{
			ID: ServerID(1),
			Coordinates: geo.Coordinates{
				Latitude:  geo.Degrees(1),
				Longitude: geo.Degrees(1),
			},
		},
		Server{
			ID: ServerID(2),
			Coordinates: geo.Coordinates{
				Latitude:  geo.Degrees(-2),
				Longitude: geo.Degrees(-2),
			},
		},
		Server{
			ID: ServerID(3),
			Coordinates: geo.Coordinates{
				Latitude:  geo.Degrees(3),
				Longitude: geo.Degrees(-3),
			},
		},
		Server{
			ID: ServerID(4),
			Coordinates: geo.Coordinates{
				Latitude:  geo.Degrees(-4),
				Longitude: geo.Degrees(4),
			},
		},
	}

	origin := geo.Coordinates{
		Latitude:  geo.Degrees(0),
		Longitude: geo.Degrees(0),
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 10; i++ {
		// Effectively a Fisher-Yates shuffle.
		//
		shuffled := make([]Server, len(expected))
		perm := r.Perm(len(shuffled))
		for i, j := range perm {
			shuffled[i] = expected[j]
		}

		allSame := true
		for i, s := range shuffled {
			if s != expected[i] {
				allSame = false
				break
			}
		}
		if allSame {
			t.Logf("Already in order on run %d", i)
		}

		_ = SortServersByDistance(shuffled, origin)
		for j, s := range shuffled {
			if s != expected[j] {
				t.Logf("Failure on run %d at index %d", i, j)
				t.Fail()
			}
		}
	}
}
