package instruments

import (
	"github.com/romansin312/metrics-exporter/tags"
)

type IInstrument interface {
	Apply(key string, n interface{}, tags ...*tags.TagModel) error
	CreateInnerInstrument(key string, description string) error
}
