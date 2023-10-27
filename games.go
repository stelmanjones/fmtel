package fmtel

import "fmt"

type Game uint

const (
	FH5 Game = iota
	FM8
	Unknown
)

func GameFromString(g string) (Game, error) {
	switch g {
	case "fh5", "FH5":
		return FH5, nil
	case "fm8", "FM8":
		return FM8, nil
	default:
		return Unknown, fmt.Errorf("expected one of (FH5,fh5,FM8,fm8). Got '%s'", g)
	}
}
