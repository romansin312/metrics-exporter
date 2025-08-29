package metrics_exporter

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
)

func convertToFloat64(n interface{}) (float64, error) {
	switch v := n.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", n)
	}
}

func convertToInt64(n interface{}) (int64, error) {
	switch v := n.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", n)
	}
}

func convertTagsToAttributes(tags []*TagModel) []attribute.KeyValue {
	if len(tags) == 0 {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, len(tags))
	for _, tag := range tags {
		attrs = append(attrs, attribute.String(tag.GetName(), tag.GetValue()))
	}

	return attrs
}
