package main

import (
	"fmt"

	"github.com/adriansr/fgnavbot/nav"
)

type scanner struct {
	line int
}

func (s *scanner) OnNext(navaid *nav.Navaid) {
	fmt.Printf("Got line %v\n", navaid)
	s.line++
}

func (s *scanner) OnError(err error) {
	fmt.Printf("Got error %s at line %d\n", err, *s)
}

func (s *scanner) OnComplete() {
	fmt.Printf("Read %d lines\n", s.line)
}

func main() {
	s := nav.Reader(&scanner{0})
	nav.Parse("nav.dat.gz", s)
}
