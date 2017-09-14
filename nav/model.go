// Package nav models the contents of flightgear's nav.dat
// i.e. X-PLANE Navigation Data File Spec version 810
package nav

import (
	"bufio"
	"compress/gzip"
	"os"
	"strconv"
	"strings"
)

type Coords struct {
	Latitude  float64
	Longitude float64
}

type NavaidType byte

const (
	NDB       NavaidType = 2
	VOR                  = 3
	LOC                  = 4 // Localizer (ILS)
	LOC_ONLY             = 5 // Localizer (No-ILS, loc only approach)
	GLS                  = 6 // Glideslope (ILS)
	OM                   = 7 // Outer marker (ILS)
	MI                   = 8 // Middle marker (ILS)
	IN                   = 9 // Inner marker (ILS)
	DME                  = 12
	DME_ONLY             = 13 // Standalone or NDB-DME
	MAX_VALUE            = 14
)

type Navaid struct {
	Type       byte
	Pos        Coords
	Elevation  int      // Feet MSL
	Frequency  int      // hundreths of MHz, or KHz for NDB
	Range      int      // Nautical miles
	Variation  float64  // Different meanings for Type
	Identifier string   // 2-3 letter code
	Extra      []string // Free-form, with parsable info if ILS component
}

type Airport struct {
	Code string
	Name string
	Pos  Coords
}

func parse(fields []string) *Navaid {
	if len(fields) < 7 {
		return nil
	}
	navType, err := strconv.Atoi(fields[0])
	if err != nil || navType <= 0 || navType >= MAX_VALUE {
		return nil
	}
	navaid := &Navaid{}
	navaid.Type = byte(navType)
	if navaid.Pos.Latitude, err = strconv.ParseFloat(fields[1], 64); err != nil {
		return nil
	}
	if navaid.Pos.Longitude, err = strconv.ParseFloat(fields[2], 64); err != nil {
		return nil
	}
	if navaid.Elevation, err = strconv.Atoi(fields[3]); err != nil {
		return nil
	}
	if navaid.Frequency, err = strconv.Atoi(fields[4]); err != nil {
		return nil
	}
	if navaid.Range, err = strconv.Atoi(fields[5]); err != nil {
		return nil
	}
	if navaid.Variation, err = strconv.ParseFloat(fields[6], 64); err != nil {
		return nil
	}
	navaid.Identifier = fields[7]
	navaid.Extra = fields[8:]
	return navaid
}

// Parse returns a...
func Parse(path string, reader Reader) {
	var handle *os.File
	var err error
	if handle, err = os.Open(path); err != nil {
		reader.OnError(err)
		return
	}
	defer handle.Close()

	var fz *gzip.Reader
	if fz, err = gzip.NewReader(handle); err != nil {
		reader.OnError(err)
		return
	}
	defer fz.Close()

	scanner := bufio.NewScanner(fz)
	for scanner.Scan() {
		if nav := parse(strings.Fields(scanner.Text())); nav != nil {
			reader.OnNext(nav)
		}
	}
	reader.OnComplete()
}
