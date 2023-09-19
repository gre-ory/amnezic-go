package util

type SqlDecoder[Model any, Row any] interface {
	DecodeRow(obj *Row) *Model
}
