package metrics_exporter

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"time"
)

const defaultUrlPath = "/v1/metrics"
const defaultSendingInterval = time.Second * 3
const defaultTimeout = time.Second * 5

type Config struct {
	endpoint              string
	urlPath               string
	headers               map[string]string
	timeout               time.Duration
	sendingInterval       time.Duration
	serviceName           string
	enableConsoleExporter bool
	useSsl                bool
}

/*
NewConfig - функция создания нового объекта конфигурации, обязательные параметры:
serviceName - название сервиса, для которого будут отправляться метрики.
endpoint - хост адреса коллектора метрик (например, localhost:1234).
*/
func NewConfig(serviceName string, endpoint string) *Config {
	newConfig := &Config{
		serviceName:     serviceName,
		endpoint:        endpoint,
		useSsl:          true,
		urlPath:         defaultUrlPath,
		timeout:         defaultTimeout,
		sendingInterval: defaultSendingInterval,
	}

	return newConfig
}

/*
WithUrlPath устанавливает путь к адресу коллектора. По умолчанию "/v1/metrics".
*/
func (c *Config) WithUrlPath(urlPath string) *Config {
	c.urlPath = urlPath
	return c
}

/*
WithHeaders устанавливает дополнительные заголовки запроса, которые будут отправлены в коллектор (например, авторизация).
*/
func (c *Config) WithHeaders(headers map[string]string) *Config {
	c.headers = headers
	return c
}

/*
WithTimeout устанавливает таймаут соединения с коллектором, по умолчанию - 5 секунд.
*/
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.timeout = timeout
	return c
}

/*
WithSendingInterval устанавливает периодичность отправки метрик в коллектор, по умолчанию - 3 секунды.
*/
func (c *Config) WithSendingInterval(sendingInterval time.Duration) *Config {
	c.sendingInterval = sendingInterval
	return c
}

/*
WithEnableConsoleExporter устанавливает, выводить ли метрики дополнительно в консоль, по умолчанию значение false.
*/
func (c *Config) WithEnableConsoleExporter(enableConsoleExporter bool) *Config {
	c.enableConsoleExporter = enableConsoleExporter
	return c
}

/*
WithUseSsl устанавливает, использовать ли SSL при соединении с коллектором, по умолчанию значение true.
*/
func (c *Config) WithUseSsl(useSsl bool) *Config {
	c.useSsl = useSsl
	return c
}

func (c *Config) build() []otlpmetrichttp.Option {
	result := make([]otlpmetrichttp.Option, 0)
	result = append(result, otlpmetrichttp.WithEndpoint(c.endpoint))
	result = append(result, otlpmetrichttp.WithTimeout(c.timeout))
	result = append(result, otlpmetrichttp.WithURLPath(c.urlPath))

	if !c.useSsl {
		result = append(result, otlpmetrichttp.WithInsecure())
	}

	if len(c.headers) > 0 {
		result = append(result, otlpmetrichttp.WithHeaders(c.headers))
	}

	return result
}
