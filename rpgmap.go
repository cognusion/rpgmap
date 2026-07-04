package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/pflag"
)

type point [2]float64
type marker struct {
	Point point
	Text  string
	Tags  []string
}

type polymarker struct {
	Points []point
	Text   string
	Tags   []string
}

type circlemarker struct {
	Point  point
	Radius float64
	Text   string
	Tags   []string
}

func newMarker(line string) (*marker, error) {
	m := marker{}

	parts := strings.Split(line, ",")
	if len(parts) < 4 {
		return nil, fmt.Errorf("number of parts must be at least 4")
	}
	m.Point = point{cast.ToFloat64(strings.TrimSpace(parts[0])), cast.ToFloat64(strings.TrimSpace(parts[1]))}
	m.Text = strings.TrimSpace(parts[2])
	m.Tags = append(m.Tags, parts[3:]...)

	return &m, nil
}

func newCirclemarker(line string) (*circlemarker, error) {
	m := circlemarker{}

	parts := strings.Split(line, ",")
	if len(parts) < 6 {
		return nil, fmt.Errorf("number of parts must be at least 6")
	}
	m.Point = point{cast.ToFloat64(strings.TrimSpace(parts[1])), cast.ToFloat64(strings.TrimSpace(parts[2]))}
	m.Radius = cast.ToFloat64(parts[3])
	m.Text = strings.TrimSpace(parts[4])
	m.Tags = append(m.Tags, parts[5:]...)

	return &m, nil
}

func newPolymarker(line string) (*polymarker, error) {
	m := polymarker{}

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

	points := make([]point, count/2)
	pi := 0
	for i := 1; i < count; i += 2 {
		points[pi] = point{cast.ToFloat64(strings.TrimSpace(parts[i])), cast.ToFloat64(strings.TrimSpace(parts[i+1]))}
		pi++
	}
	m.Points = points
	m.Text = strings.TrimSpace(parts[count+1])
	m.Tags = append(m.Tags, parts[count+2:]...)

	return &m, nil
}

func main() {

	tags := make(map[string][]string)

	var (
		conf = pflag.StringP("config", "c", "", "Config file to read")
	)
	pflag.Parse()

	f, err := os.Open(*conf)
	if err != nil {
		panic(err) // TODO
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var (
		line         string
		commentBlock bool
	)
	for scanner.Scan() {
		// We always trimp starting and ending whitespace
		line = strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}
		// Skip comments
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "/*") {
			commentBlock = true
			continue
		}
		if strings.HasPrefix(line, "*/") {
			commentBlock = false
			continue
		}
		if commentBlock {
			continue
		}

		// Process it!
		if strings.HasPrefix(line, "c") {
			// circle
			m, err := newCirclemarker(line)
			if err != nil {
				panic(err)
			}

			mLine := fmt.Sprintf("L.circle([%f,%f],{ radius: %f, color: 'red', fillColor: '#f03', fillOpacity: 0.2}).bindPopup(\"%s\")", m.Point[0], m.Point[1], m.Radius, m.Text)
			for _, t := range m.Tags {
				t = strings.TrimSpace(t)
				tags[t] = append(tags[t], mLine)
			}

		} else if strings.HasPrefix(line, "p") {
			// Polygon
			m, err := newPolymarker(line)
			if err != nil {
				panic(err)
			}

			mLine := "L.polygon(["
			for i, p := range m.Points {
				mLine += fmt.Sprintf("  [%f,%f]", p[0], p[1])
				if i+1 < len(m.Points) {
					mLine += ","
				}
			}
			mLine += fmt.Sprintf("],{ color: 'red', fillColor: '#f03', fillOpacity: 0.2}).bindPopup(\"%s\")", m.Text)

			for _, t := range m.Tags {
				t = strings.TrimSpace(t)
				tags[t] = append(tags[t], mLine)
			}

		} else {
			// Marker
			m, err := newMarker(line)
			if err != nil {
				panic(err)
			}

			mLine := fmt.Sprintf("L.marker([%f,%f]).bindPopup(\"%s\")", m.Point[0], m.Point[1], m.Text)
			for _, t := range m.Tags {
				t = strings.TrimSpace(t)
				tags[t] = append(tags[t], mLine)
			}
		}

	}
	if err := scanner.Err(); err != nil {
		panic(err) //TODO
	}

	// Dump
	fmt.Printf("var layerControl = L.control.layers().addTo(map)\n")
	sortedTags := slices.Sorted(maps.Keys(tags))
	for _, t := range sortedTags {
		lines := tags[t]
		fmt.Printf("var %s = L.layerGroup([\n", t)
		for i, l := range lines {
			fmt.Printf("\t%s", l)
			if i+1 < len(lines) {
				fmt.Println(",")
			}
		}
		fmt.Println("])")
		fmt.Printf("layerControl.addOverlay(%s, \"%s\");\n", t, strings.Title(t))
	}

}
