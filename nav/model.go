package nav

type Coords struct {
	Latitude float64
	Longitude float64
}

type NavaidType byte

const (
	NDB NavaidType = iota
	VOR
	LOC
	GLS
	DME

)
type Navaid struct {
	Type byte
}