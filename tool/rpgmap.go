package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/cognusion/rpgmap"
	"github.com/spf13/pflag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {

	var (
		conf = pflag.StringP("config", "c", "", "Config file to read")
	)
	pflag.Parse()

	f, err := os.Open(*conf)
	if err != nil {
		dief("Error opening config file '%s': %s\n", *conf, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var (
		tags         = make(map[string][]string)
		icons        = make(map[string]rpgmap.Icon)
		line         string
		commentBlock bool
		altMaps      = make([]rpgmap.Map, 0)
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
			m   rpgmap.TagStringer
			err error
		)

		// Build based on the line type
		if strings.HasPrefix(line, "a") {
			// Altmap!
			a, ae := rpgmap.NewMap(line)
			if ae != nil {
				die(ae)
			}
			altMaps = append(altMaps, *a)
			continue // we're done with this line

		} else if strings.HasPrefix(line, "i") {
			// icon!!
			i, ie := rpgmap.NewIcon(line)
			if ie != nil {
				die(ie)
			}
			icons[i.Tag] = *i
			continue // we don't want to continue on

		} else if strings.HasPrefix(line, "c") {
			// circle
			m, err = rpgmap.NewCircleMarker(line)

		} else if strings.HasPrefix(line, "p") {
			// Polygon
			m, err = rpgmap.NewPolyMarker(line)

		} else {
			// Point
			m, err = rpgmap.NewPointMarker(line, icons)

		}

		// Handle err
		if err != nil {
			die(err)
		}

		// Hook it into the tag map
		tags = addLineToTags(tags, m.String(), m.Tags())

	}
	if err := scanner.Err(); err != nil {
		die(err)
	}

	// Dump
	fmt.Printf("var layerControl = L.control.layers().addTo(map)\n")

	for _, altmap := range altMaps {
		fmt.Printf("%s\n", altmap.String())
	}

	for tag, icon := range icons {
		fmt.Printf("var %sIcon = %s;\n", tag, icon.String())
	}

	sortedTags := slices.Sorted(maps.Keys(tags))
	for _, t := range sortedTags {
		lines := tags[t]
		tag := rpgmap.CleanTag(t)
		fmt.Printf("var %s = L.layerGroup([\n", tag)
		for i, l := range lines {
			fmt.Printf("\t%s", l)
			if i+1 < len(lines) {
				fmt.Println(",")
			}
		}
		fmt.Println("])")
		fmt.Printf("layerControl.addOverlay(%s, \"%s\");\n", tag, title(t))
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

// strings.Title is deprecated, so we have recreated it, because that's better.
func title(in string) string {
	engCaser := cases.Title(language.English)
	return engCaser.String(in)
}

func die(e error) {
	dief("%s\n", e)
}

func dief(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
