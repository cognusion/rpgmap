package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/spf13/pflag"
)

func main() {

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
		tags         = make(map[string][]string)
		line         string
		commentBlock bool
	)
	for scanner.Scan() {
		// We always trim starting and ending whitespace
		line = strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}
		// Skip comments
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		} else if strings.HasPrefix(line, "/*") {
			commentBlock = true
			continue
		} else if strings.HasPrefix(line, "*/") {
			commentBlock = false
			continue
		}
		if commentBlock {
			continue
		}

		// Process it!
		if strings.HasPrefix(line, "c") {
			// circle
			m, err := NewCircleMarker(line)
			if err != nil {
				panic(err)
			}

			mLine := fmt.Sprintf("L.circle(%s,{ radius: %f, color: 'red', fillColor: '#f03', fillOpacity: 0.2}).bindPopup(\"%s\")", m.Point.String(), m.Radius, m.Text)
			for _, t := range m.Tags {
				t = strings.TrimSpace(t)
				tags[t] = append(tags[t], mLine)
			}

		} else if strings.HasPrefix(line, "p") {
			// Polygon
			m, err := NewPolyMarker(line)
			if err != nil {
				panic(err)
			}

			mLine := "L.polygon(["
			for i, p := range m.Points {
				mLine += p.String()
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
			m, err := NewPointMarker(line)
			if err != nil {
				panic(err)
			}

			mLine := fmt.Sprintf("L.marker(%s).bindPopup(\"%s\")", m.Point.String(), m.Text)
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
