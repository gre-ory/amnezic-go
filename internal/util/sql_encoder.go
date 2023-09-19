package util

type SqlEncoder[Model any, Row any] interface {
	EncodeRow(obj *Model) *Row
}
