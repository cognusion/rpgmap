package main

import (
	"bufio"
	"fmt"
	"io"
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
		litLines     = make([]string, 0)
		lineCount    int64
	)

	// PREPROCESS for icons
	for scanner.Scan() {
		// We always trim starting and ending whitespace
		line = strings.TrimSpace(scanner.Text())
		lineCount++

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

		if strings.HasPrefix(line, "i") {
			// icon!!
			i, ie := rpgmap.NewIcon(line)
			if ie != nil {
				dief("preprocess error on line %d: %v\n", lineCount, ie)
			}
			icons[i.Tag] = *i
		}
	}
	if err = scanner.Err(); err != nil {
		die(err)
	}

	// Reset the file to process again
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		die(err)
	}
	lineCount = 0 // reset
	scanner = bufio.NewScanner(f)

	// second pass for effect
	for scanner.Scan() {
		// We always trim starting and ending whitespace
		line = strings.TrimSpace(scanner.Text())
		lineCount++

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
			m      rpgmap.TagStringer
			newErr error
		)

		// Build based on the line type
		if l, cut := strings.CutPrefix(line, "!"); cut {
			// Literal line
			litLines = append(litLines, l)
			continue // we're done with this line

		} else if strings.HasPrefix(line, "a") {
			// Altmap!
			a, ae := rpgmap.NewMap(line)
			if ae != nil {
				dief("error on line %d: %v\n", lineCount, ae)
			}
			altMaps = append(altMaps, *a)
			continue // we're done with this line

		} else if strings.HasPrefix(line, "i") {
			// icon, which we preprocessed. Skip!
			continue

		} else if strings.HasPrefix(line, "c") {
			// circle
			m, newErr = rpgmap.NewCircleMarker(line)

		} else if strings.HasPrefix(line, "p") {
			// Polygon
			m, newErr = rpgmap.NewPolyMarker(line)

		} else {
			// Point
			m, newErr = rpgmap.NewPointMarker(line, icons)

		}

		// Handle err
		if newErr != nil {
			dief("error on line %d: %v\n", lineCount, newErr)
		}

		// Hook it into the tag map
		tags = addLineToTags(tags, m.String(), m.Tags())

	}
	if err = scanner.Err(); err != nil {
		die(err)
	}

	// Dump
	fmt.Printf("var layerControl = L.control.layers().addTo(map)\n\n")

	for _, altmap := range altMaps {
		fmt.Printf("%s\n", altmap.String())
	}

	sortedIcons := slices.Sorted(maps.Keys(icons))
	for _, tag := range sortedIcons {
		icon := icons[tag]
		fmt.Printf("var %sIcon = %s;\n", tag, icon.String())
	}
	fmt.Println()

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
		fmt.Printf("layerControl.addOverlay(%s, \"%s\");\n\n", tag, title(t))
	}

	fmt.Println()
	for _, l := range litLines {
		fmt.Println(l)
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
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
