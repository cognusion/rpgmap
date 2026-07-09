package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
		case bool:
			fmt.Fprintf(&l, "%s: %t", opt.Option, t)
		case string:
			if ct, bare := strings.CutPrefix(t, "BARE"); bare {
				// Don't quote it
				fmt.Fprintf(&l, "%s: %s", opt.Option, ct)
			} else {
				// Quote it
				fmt.Fprintf(&l, "%s: '%s'", opt.Option, t)
			}
		case int64:
			fmt.Fprintf(&l, "%s: %d", opt.Option, t)
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

// ParseOptString treats a string as if it defines a key:value pair,
// and returns the appropriate Opt or an error
func ParseOptString(ostring string) (Opt, error) {
	var (
		o Opt
	)

	parts := strings.SplitN(ostring, ":", 2)
	if len(parts) != 2 {
		return o, fmt.Errorf("option string invalid")
	}

	o.Option = parts[0]

	if i, err := strconv.ParseInt(parts[1], 0, 64); err == nil {
		o.Value = i
	} else if f, err := strconv.ParseFloat(parts[1], 64); err == nil {
		o.Value = f
	} else if b, err := strconv.ParseBool(parts[1]); err == nil {
		o.Value = b
	} else {
		// Assume it is a string
		o.Value = parts[1]
	}

	return o, nil
}
