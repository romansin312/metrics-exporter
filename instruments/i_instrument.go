package instruments

import (
	"github.com/romansin312/metrics-exporter/tags"
)

type IInstrument interface {
	Apply(key string, n interface{}, tags ...*tags.TagModel) error
	Create(key string, description string) error
}
