package rpgmap

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

// Point is a float64 tuple referring to an X/Y coordinate.
type Point [2]float64

// String formats the Point and return the string representation encased in hard braces.
func (p Point) String() string {
	return fmt.Sprintf("[%f,%f]", p[0], p[1])
}

// TagsToOpts takes a list of tags and returns a list of tags and Opts if any.
func TagsToOpts(tags ...string) ([]string, []Opt) {
	var (
		os = make([]Opt, 0)
		ts = make([]string, 0)
	)

	for _, t := range tags {
		if strings.Contains(t, ":") {
			o, _ := ParseOptString(t) // We are deliberately ignoring errs because we have mitigated the only returned error already
			os = append(os, o)
		} else {
			ts = append(ts, t)
		}
	}
	return ts, os
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

// Map is a type to describe an overlay map
type Map struct {
	Name string
	URI  string
	Options
}

func (m *Map) String() string {
	var mLine strings.Builder
	cname := CleanTag(m.Name)
	fmt.Fprintf(&mLine, "var %s = L.tileLayer('%s'", cname, m.URI)
	if !m.Options.IsEmpty() {
		fmt.Fprintf(&mLine, ",%s", m.Options.String())
	}
	fmt.Fprintf(&mLine, ")\nlayerControl.addOverlay(%s, \"%s\");", cname, m.Name)
	return mLine.String()
}

// Icon is a type to describe a tag-associated PointMarker icon
type Icon struct {
	Tag        string
	URI        string
	Size       Point
	IconAnchor Point
}

func (i *Icon) String() string {
	return fmt.Sprintf("L.icon({ iconUrl: '%s', iconSize: %s, iconAnchor: %s })", i.URI, i.Size, i.IconAnchor)
	// We can always get away with printing iconanchor because the default is 0,0
}

// PointMarker is a Marker referred to by a single Point.
type PointMarker struct {
	Point Point
	Marker
	IconTag string
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
	fmt.Fprintf(&mLine, ").bindPopup(\"%s\")", m.Text)
	return mLine.String()
}

// NewMap takes a line string and return a Map or an error.
func NewMap(line string) (*Map, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 3 {
		return nil, fmt.Errorf("not enough parameters in line")
	}
	m := Map{
		Name: parts[1],
		URI:  parts[2],
	}
	for _, p := range parts[3:] {
		o, e := ParseOptString(p)
		if e != nil {
			return nil, fmt.Errorf("error parsing an option: %w", e)
		}
		m.Options.Add(o)
	}
	return &m, nil
}

// NewIcon takes a line string and returns an Icon or an error.
func NewIcon(line string) (*Icon, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 5 {
		return nil, fmt.Errorf("not enough parameters in line")
	}

	i := Icon{
		Tag:        CleanTag(parts[1]),
		URI:        parts[2],
		Size:       Point{cast.ToFloat64(strings.TrimSpace(parts[3])), cast.ToFloat64(strings.TrimSpace(parts[4]))},
		IconAnchor: Point{0, 0}, //default origin
	}
	if len(parts) == 7 {
		// anchor point specified
		i.IconAnchor = Point{cast.ToFloat64(strings.TrimSpace(parts[5])), cast.ToFloat64(strings.TrimSpace(parts[6]))}
	}

	return &i, nil
}

// NewPointMarker takes a line string and returns a PointMarker or an error.
func NewPointMarker(line string, icons map[string]Icon) (*PointMarker, error) {
	m := PointMarker{}

	parts := strings.Split(line, ",")
	if len(parts) < 4 {
		return nil, fmt.Errorf("number of parts must be at least 4")
	}
	m.Point = Point{cast.ToFloat64(strings.TrimSpace(parts[0])), cast.ToFloat64(strings.TrimSpace(parts[1]))}
	m.Text = strings.TrimSpace(parts[2])

	tags, opts := TagsToOpts(parts[3:]...)

	if len(tags) > 0 {
		m.AddTags(tags...)

		if icons != nil {
			tag := CleanTag(tags[0])
			if _, ok := icons[tag]; ok {
				m.Options.Add(Opt{"icon", fmt.Sprintf("BARE%sIcon", tag)})
			}
		}
	}
	if len(opts) > 0 {
		m.Options.Add(opts...)
	}

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
	m.Options.Add(Opt{"radius", m.Radius})

	tags, opts := TagsToOpts(parts[5:]...)

	if len(tags) > 0 {
		m.AddTags(tags...)
	}
	if len(opts) > 0 {
		m.Options.Add(opts...)
	}

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

	tags, opts := TagsToOpts(parts[count+2:]...)

	if len(tags) > 0 {
		m.AddTags(tags...)
	}
	if len(opts) > 0 {
		m.Options.Add(opts...)
	}

	return &m, nil
}

// CleanTag removes all non-letters from a tag
func CleanTag(tag string) string {
	reg := regexp.MustCompile(`[^\p{L}]+`)
	return reg.ReplaceAllString(tag, "")
}
