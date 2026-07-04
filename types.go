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

// Options is a slice of Opts, for all Markers to inherit
type Options []Opt

// Opt is a {string,any} tuple
type Opt struct {
	Option string
	Value  any
}

// Add appends slices of Opt
func (o *Options) Add(options ...Opt) {
	*o = append(*o, options...)
}

// IsEmpty return true if Options has no entries.
func (o *Options) IsEmpty() bool {
	return len(*o) == 0
}

// String returns the stringified list of Opt.
func (o *Options) String() string {
	var l strings.Builder
	l.WriteString("{")

	c := 0
	for _, opt := range *o {
		if c > 0 {
			l.WriteString(",")
		}
		c++

		switch t := opt.Value.(type) {
		case string:
			fmt.Fprintf(&l, "%s: '%s'", opt.Option, t)
		case float64:
			fmt.Fprintf(&l, "%s: %f", opt.Option, t)
		default:
			// probably should error, but...
			fmt.Fprintf(&l, "%s: '%v", opt.Option, t)
		}
	}
	l.WriteString("}")
	return l.String()
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
	Options
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
	var mLine string
	mLine = m.Point.String()
	if !m.Options.IsEmpty() {
		mLine += "," + m.Options.String()
	}
	return fmt.Sprintf("L.marker(%s).bindPopup(\"%s\")", mLine, m.Text)
}

// PolyMarker is a Marker referred to by a list of Point creating a polygon.
type PolyMarker struct {
	Points []Point
	Marker
}

func (m PolyMarker) String() string {
	var mLine strings.Builder
	mLine.WriteString("L.polygon([")
	for i, p := range m.Points {
		mLine.WriteString(p.String())
		if i+1 < len(m.Points) {
			mLine.WriteString(",")
		}
	}
	mLine.WriteString("]")
	if !m.Options.IsEmpty() {
		fmt.Fprintf(&mLine, ",%s", m.Options.String())
	}
	fmt.Fprintf(&mLine, ").bindPopup(\"%s\")", m.Text)

	return mLine.String()
}

// CircleMarker is a Marker referred to by a Point centroid and a radius from that point.
type CircleMarker struct {
	Point  Point
	Radius float64
	Marker
}

func (m CircleMarker) String() string {
	var mLine strings.Builder
	fmt.Fprintf(&mLine, "L.circle(%s", m.Point.String())
	if !m.Options.IsEmpty() {
		fmt.Fprintf(&mLine, ",%s", m.Options.String())
	}
	fmt.Fprintf(&mLine, ".bindPopup(\"%s\")", m.Text)
	return mLine.String()
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

	m.Options.Add(Opt{"color", "red"}, Opt{"fillColor", "#f03"}, Opt{"radius", m.Radius}, Opt{"fillOpacity", 0.2})
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

	m.Options.Add(Opt{"color", "red"}, Opt{"fillColor", "#f03"}, Opt{"fillOpacity", 0.2})
	return &m, nil
}
