package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

// Point is a float64 tuple referring to an X/Y coordinate.
type Point [2]float64

// String formats the Point and return the string representation encased in hard braces.
func (p Point) String() string {
	return fmt.Sprintf("[%f,%f]", p[0], p[1])
}

// TagStringer is a nonsense interface that defines the only two things we really need
// from these types in order to generate the text we need.
type TagStringer interface {
	String() string
	Tags() []string
}

// Marker is a type that has a description and a list of tags.
type Marker struct {
	Text string
	tags []string
}

// AddTags takes an arbitrary number of strings and appends them to
// the internal tag list
func (m *Marker) AddTags(tags ...string) {
	m.tags = append(m.tags, tags...)
}

// Tags returns the internal tag list
func (m *Marker) Tags() []string {
	return m.tags
}

// PointMarker is a Marker referred to by a single Point.
type PointMarker struct {
	Point Point
	Marker
}

func (m PointMarker) String() string {
	return fmt.Sprintf("L.marker(%s).bindPopup(\"%s\")", m.Point.String(), m.Text)
}

// PolyMarker is a Marker referred to by a list of Point creating a polygon.
type PolyMarker struct {
	Points []Point
	Marker
}

func (m PolyMarker) String() string {
	mLine := "L.polygon(["
	for i, p := range m.Points {
		mLine += p.String()
		if i+1 < len(m.Points) {
			mLine += ","
		}
	}
	mLine += fmt.Sprintf("],{ color: 'red', fillColor: '#f03', fillOpacity: 0.2}).bindPopup(\"%s\")", m.Text)

	return mLine
}

// CircleMarker is a Marker referred to by a Point centroid and a radius from that point.
type CircleMarker struct {
	Point  Point
	Radius float64
	Marker
}

func (m CircleMarker) String() string {
	return fmt.Sprintf("L.circle(%s,{ radius: %f, color: 'red', fillColor: '#f03', fillOpacity: 0.2}).bindPopup(\"%s\")", m.Point.String(), m.Radius, m.Text)

}

// NewPointMarker takes a line string and returns a PointMarker or an error.
func NewPointMarker(line string) (*PointMarker, error) {
	m := PointMarker{}

	parts := strings.Split(line, ",")
	if len(parts) < 4 {
		return nil, fmt.Errorf("number of parts must be at least 4")
	}
	m.Point = Point{cast.ToFloat64(strings.TrimSpace(parts[0])), cast.ToFloat64(strings.TrimSpace(parts[1]))}
	m.Text = strings.TrimSpace(parts[2])
	m.AddTags(parts[3:]...)

	return &m, nil
}

// NewCircleMarker takes a line string and returns a CircleMarker or an error.
func NewCircleMarker(line string) (*CircleMarker, error) {
	m := CircleMarker{}

	parts := strings.Split(line, ",")
	if len(parts) < 6 {
		return nil, fmt.Errorf("number of parts must be at least 6")
	}
	m.Point = Point{cast.ToFloat64(strings.TrimSpace(parts[1])), cast.ToFloat64(strings.TrimSpace(parts[2]))}
	m.Radius = cast.ToFloat64(parts[3])
	m.Text = strings.TrimSpace(parts[4])
	m.AddTags(parts[5:]...)

	return &m, nil
}

// NewPolyMarker takes a line string and returns a PolyMarker or an error.
func NewPolyMarker(line string) (*PolyMarker, error) {
	m := PolyMarker{}

	parts := strings.Split(line, ",")
	if len(parts) < 5 {
		return nil, fmt.Errorf("number of parts must be at least 5")
	}

	// divine the number of points
	var count int
	for _, p := range parts[1:] {
		if _, e := cast.ToFloat64E(strings.TrimSpace(p)); e != nil {
			break
		}
		count++
	}
	if count%2 != 0 {
		return nil, fmt.Errorf("polygon points must have an even number, not %d", count)
	}

	points := make([]Point, count/2)
	pi := 0
	for i := 1; i < count; i += 2 {
		points[pi] = Point{cast.ToFloat64(strings.TrimSpace(parts[i])), cast.ToFloat64(strings.TrimSpace(parts[i+1]))}
		pi++
	}
	m.Points = points
	m.Text = strings.TrimSpace(parts[count+1])
	m.AddTags(parts[count+2:]...)

	return &m, nil
}
