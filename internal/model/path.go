package model

import (
	"path/filepath"
)

// //////////////////////////////////////////////////
// path validator

type PathValidator func(path string) error

// //////////////////////////////////////////////////
// helper

func Path(elems ...string) Url {
	return Url(filepath.Join(elems...))
}
