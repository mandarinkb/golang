package config

import "time"

type Root struct {
	MariaDb MariaDB `mapstructure:"mariadb"`
	Elastic Elastic `mapstructure:"elastic"`
}

type MariaDB struct {
	MariaDBPassword string `mapstructure:"maria_password"`
	DriverName      string `mapstructure:"driver_name"`
	DataSourceName  string `mapstructure:"data_source_name"`
}
type Elastic struct {
	HTTPTransport  HTTPTransport `mapstructure:"http_transport"`
	ServerWithPort string        `mapstructure:"server_with_port"`
	Index          string        `mapstructure:"index"`
	Type           string        `mapstructure:"type"`
}

type HTTPTransport struct {
	Timeout               time.Duration `mapstructure:"timeout"`
	SkipVerifyTLS         bool          `mapstructure:"skip_verify_tls"`
	DialTimeout           time.Duration `mapstructure:"dial_timeout"`
	DialKeepAlive         time.Duration `mapstructure:"dial_keep_alive"`
	MaxIdleConns          int           `mapstructure:"max_idle_conns"`
	MaxIdleConnsPerHost   int           `mapstructure:"max_idle_conns_per_host"`
	IdleConnTimeout       time.Duration `mapstructure:"idle_conn_timeout"`
	TLSHandshakeTimeout   time.Duration `mapstructure:"tls_handshake_timeout"`
	ResponseHeaderTimeout time.Duration `mapstructure:"respose_header_timeout"`
	ExpectContinueTimeout time.Duration `mapstructure:"expect_continue_timeout"`
}
