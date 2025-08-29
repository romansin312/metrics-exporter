package metrics_exporter

type IInstrument interface {
	Apply(key string, n interface{}, tags ...*TagModel) error
	Create(key string, description string) error
}
