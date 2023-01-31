package model

import "fmt"

// //////////////////////////////////////////////////
// error

var (
	ErrNotImplemented = fmt.Errorf("not implemented")
	ErrGameNotFound   = fmt.Errorf("game not found")
)
