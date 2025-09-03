package metrics_exporter

import (
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"os"
	"strings"
	"time"
)

const defaultUrlPath = "/v1/metrics"
const defaultSendingInterval = time.Second * 3
const defaultTimeout = time.Second * 5
const defaultScope = "default"

type Config struct {
	host                  string
	urlPath               string
	headers               map[string]string
	timeout               time.Duration
	sendingInterval       time.Duration
	serviceName           string
	enableConsoleExporter bool
	useSsl                bool
	scope                 string
}

type ConfigOption func(*Config)

/*
NewConfig - функция создания нового объекта конфигурации.

Если название сервиса и хост не будут переданы как ConfigOption иди прописаны в переменных окружения, метод вернёт ошибку.
*/
func NewConfig(opts ...ConfigOption) (*Config, error) {
	newConfig := &Config{
		serviceName:     "",
		host:            "",
		useSsl:          true,
		urlPath:         defaultUrlPath,
		timeout:         defaultTimeout,
		sendingInterval: defaultSendingInterval,
		scope:           defaultScope,
	}

	newConfig.loadFromEnv()

	for _, opt := range opts {
		opt(newConfig)
	}

	if newConfig.serviceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	if newConfig.host == "" {
		return nil, fmt.Errorf("host is required")
	}

	return newConfig, nil
}

/*
WithUrlPath устанавливает путь к адресу коллектора. По умолчанию "/v1/metrics".
*/
func WithUrlPath(urlPath string) ConfigOption {
	return func(c *Config) {
		if urlPath != "" {
			c.urlPath = urlPath
		}
	}
}

/*
WithHeaders устанавливает дополнительные заголовки запроса, которые будут отправлены в коллектор (например, авторизация).
*/
func WithHeaders(headers map[string]string) ConfigOption {
	return func(c *Config) {
		c.headers = headers
	}
}

/*
WithTimeout устанавливает таймаут соединения с коллектором, по умолчанию - 5 секунд.
*/
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) {
		if timeout > 0 {
			c.timeout = timeout
		}
	}
}

/*
WithSendingInterval устанавливает периодичность отправки метрик в коллектор, по умолчанию - 3 секунды.
*/
func WithSendingInterval(sendingInterval time.Duration) ConfigOption {
	return func(c *Config) {
		if sendingInterval > 0 {
			c.sendingInterval = sendingInterval
		}
	}
}

/*
WithEnableConsoleExporter устанавливает, выводить ли метрики дополнительно в консоль, по умолчанию значение false.
*/
func WithEnableConsoleExporter(enableConsoleExporter bool) ConfigOption {
	return func(c *Config) {
		c.enableConsoleExporter = enableConsoleExporter
	}
}

/*
WithUseSsl устанавливает, использовать ли SSL при соединении с коллектором, по умолчанию значение true.
*/
func WithUseSsl(useSsl bool) ConfigOption {
	return func(c *Config) {
		c.useSsl = useSsl
	}
}

/*
 */
func WithScope(scope string) ConfigOption {
	return func(c *Config) {
		c.scope = scope
	}
}

/*
WithServiceName устанавливает название сервиса, для которого будут отправляться метрики
*/
func WithServiceName(serviceName string) ConfigOption {
	return func(c *Config) {
		c.serviceName = serviceName
	}
}

/*
WithHost - хост адреса коллектора метрик (например, localhost:1234).
*/
func WithHost(host string) ConfigOption {
	return func(c *Config) {
		c.host = host
	}
}

func (c *Config) build() []otlpmetrichttp.Option {
	result := make([]otlpmetrichttp.Option, 0)
	result = append(result, otlpmetrichttp.WithEndpoint(c.host))
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

func (c *Config) loadFromEnv() {
	if envEndpoint := os.Getenv("METRICS_ENDPOINT"); envEndpoint != "" {
		c.host = envEndpoint
	}

	if envUrlPath := os.Getenv("METRICS_URL_PATH"); envUrlPath != "" {
		c.urlPath = envUrlPath
	}

	if envTimeout := os.Getenv("METRICS_TIMEOUT"); envTimeout != "" {
		if timeout, err := time.ParseDuration(envTimeout); err == nil {
			c.timeout = timeout
		}
	}

	if envInterval := os.Getenv("METRICS_SENDING_INTERVAL"); envInterval != "" {
		if interval, err := time.ParseDuration(envInterval); err == nil {
			c.sendingInterval = interval
		}
	}

	if envSSL := os.Getenv("METRICS_USE_SSL"); envSSL != "" {
		c.useSsl = strings.ToLower(envSSL) == "true" || envSSL == "1"
	}

	if envConsole := os.Getenv("METRICS_ENABLE_CONSOLE"); envConsole != "" {
		c.enableConsoleExporter = strings.ToLower(envConsole) == "true" || envConsole == "1"
	}

	if envServiceName := os.Getenv("METRICS_SERVICE_NAME"); envServiceName != "" {
		c.serviceName = envServiceName
	}

	if envHeaders := os.Getenv("METRICS_HEADERS"); envHeaders != "" {
		headers := make(map[string]string)
		pairs := strings.Split(envHeaders, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
		c.headers = headers
	}
}
