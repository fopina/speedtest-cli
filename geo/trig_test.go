package geo

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestKilometers_Less(t *testing.T) {
	if Kilometers(5).Less(Kilometers(4)) {
		t.Fail()
	}
}

func ExampleKilometers_String() {
	fmt.Println(Kilometers(5))
	// Output:
	// 5.00 km
}

// Distance function should be a metric and so the triangle inequality must
// hold:
//
// for any points on the Earth (a, b, c), d(a, b) <= d(a, c) + d(b, c)
//
func TestCoordinates_DistanceTo_TriangleInequality(t *testing.T) {
	const iterations = 10

	r := rand.New(rand.NewSource(time.Now().Unix()))
	gen := func() Coordinates {
		return Coordinates{
			Latitude:  Degrees(r.Float64() * 360),
			Longitude: Degrees(r.Float64() * 360),
		}
	}

	for i := 0; i < iterations; i++ {
		a, b, c := gen(), gen(), gen()
		if a.DistanceTo(b) > a.DistanceTo(c)+b.DistanceTo(c) {
			t.Fail()
			t.Logf("Triangle inequality did not hold for a=%v, b=%v, c=%v", a, b, c)
		}
	}
}

func TestDegrees_ToRadians(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	gen := func() Degrees {
		return Degrees(r.Float64() * 360)
	}

	for i := 0; i < 10; i++ {
		if gen().ToRadians() > 2*math.Pi {
			t.Fail()
		}
	}
}
