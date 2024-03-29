package model

import "strings"

// //////////////////////////////////////////////////
// source

type Source string

var (
	Source_Legacy Source = "legacy"
	Source_Decade Source = "decade"
	Source_Genre  Source = "genre"
	Source_Store  Source = "store"
	Source_Deezer Source = "deezer"
)

func (o Source) IsStore() bool {
	return o == Source_Store
}

func (o Source) IsDeezer() bool {
	return o == Source_Deezer
}

func ToSource(value string) Source {
	value = strings.Trim(value, " ")
	value = strings.ToLower(value)
	switch value {
	case string(Source_Legacy):
		return Source_Legacy
	case string(Source_Decade):
		return Source_Decade
	case string(Source_Genre):
		return Source_Genre
	case string(Source_Store):
		return Source_Store
	case string(Source_Deezer):
		return Source_Deezer
	default:
		return ""
	}
}

func (s Source) String() string {
	return string(s)
}
