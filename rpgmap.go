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
		var (
			m   TagStringer
			err error
		)

		// Build based on the line type
		if strings.HasPrefix(line, "c") {
			// circle
			m, err = NewCircleMarker(line)

		} else if strings.HasPrefix(line, "p") {
			// Polygon
			m, err = NewPolyMarker(line)

		} else {
			// Point
			m, err = NewPointMarker(line)
		}

		// Handle err
		if err != nil {
			panic(err)
		}

		// Hook it into the tag map
		tags = addLineToTags(tags, m.String(), m.Tags())

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

// addLineToTags iterates over the tags, and adds the line to each valid entry in the tagMap
func addLineToTags(tagMap map[string][]string, line string, tags []string) map[string][]string {
	for _, t := range tags {
		t = strings.TrimSpace(t)
		tagMap[t] = append(tagMap[t], line)
	}

	return tagMap
}
