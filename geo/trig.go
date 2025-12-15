/*
A super tiny support library for go.jonnrb.io/speedtest that calculates
distances (as the crow flies) between surface points on Earth.
*/
package geo

import (
	"fmt"
	"math"
)

type Degrees float32
type Radians float64
type Kilometers float64

type Coordinates struct {
	Latitude  Degrees
	Longitude Degrees
}

const RadiusOfEarth = Kilometers(6371)

func (k Kilometers) Less(other Kilometers) bool {
	return float64(k) < float64(other)
}

func (k Kilometers) String() string {
	return fmt.Sprintf("%.2f km", float64(k))
}

func (d Degrees) ToRadians() Radians {
	return Radians(float64(d) * math.Pi / 180)
}

func (theta Radians) Sin() float64 {
	return math.Sin(float64(theta))
}

func (theta Radians) Cos() float64 {
	return math.Cos(float64(theta))
}

// From https://www.movable-type.co.uk/scripts/latlong.html
func (org Coordinates) DistanceTo(dest Coordinates) Kilometers {
	dlat := (dest.Latitude - org.Latitude).ToRadians()
	dlon := (dest.Longitude - org.Longitude).ToRadians()

	a := (dlat/2).Sin()*(dlat/2).Sin() +
		org.Latitude.ToRadians().Cos()*dest.Latitude.ToRadians().Cos()*
			(dlon/2).Sin()*(dlon/2).Sin()

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return RadiusOfEarth * Kilometers(c)
}
