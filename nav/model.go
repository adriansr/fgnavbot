// Package nav models the contents of flightgear's nav.dat.gz and apt.dat.gz
// i.e. X-PLANE Navigation Data File Spec version 810
// and X-PLANE Airport Data File Spec version 1000
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
	// No position -- Pos       Coords
	Elevation int
}

type RunwayEnd struct {
	Code         string
	Pos          Coords
	DispThMeters float64 // displaced threshold, in meters!
	// No length -- Length int
}

type Runway struct {
	Width float64
	End   [2]RunwayEnd
}

// To signal end of data
type Terminator int

func parseNavaid(fields []string) interface{} {
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

func parseAirport(fields []string) interface{} {
	if len(fields) < 1 {
		return nil
	}
	entryCode, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil
	}

	if entryCode == 1 && len(fields) > 5 {
		airport := &Airport{}
		if airport.Elevation, err = strconv.Atoi(fields[1]); err != nil {
			return nil
		}
		airport.Code = fields[4]
		airport.Name = strings.Join(fields[5:], " ")
		// TODO: missing airport.Pos // deduce from runway ends?
		return airport
	}
	if entryCode == 100 && len(fields) == 26 {
		runway := &Runway{}
		if runway.Width, err = strconv.ParseFloat(fields[1], 64); err != nil {
			return nil
		}
		for i := 0; i < 2; i++ {
			pos := 8 + 9*i
			end := &runway.End[i]
			end.Code = fields[pos]
			if end.Pos.Latitude, err = strconv.ParseFloat(fields[pos+1], 64); err != nil {
				return err
			}
			if end.Pos.Longitude, err = strconv.ParseFloat(fields[pos+2], 64); err != nil {
				return err
			}
			if end.DispThMeters, err = strconv.ParseFloat(fields[pos+3], 64); err != nil {
				return err
			}
		}
		return runway
	}
	return nil
}

func gzTextParser(path string, parser func([]string) interface{}, out chan interface{}) {
	var handle *os.File
	var err error
	if handle, err = os.Open(path); err != nil {
		out <- err
		return
	}
	defer handle.Close()

	var fz *gzip.Reader
	if fz, err = gzip.NewReader(handle); err != nil {
		out <- err
		return
	}
	defer fz.Close()

	scanner := bufio.NewScanner(fz)
	for scanner.Scan() {
		if nav := parser(strings.Fields(scanner.Text())); nav != nil {
			out <- nav
		}
	}
	out <- Terminator(0)
}

// Parse returns a...
func ReadNavaids(path string, out chan interface{}) {
	gzTextParser(path, parseNavaid, out)
}

// Parse returns a...
func ReadAirports(path string, out chan interface{}) {
	gzTextParser(path, parseAirport, out)
}
