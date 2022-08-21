package config

type Root struct {
	MariaDb MariaDB `mapstructure:"mariadb"`
}

type MariaDB struct {
	MariaDBPassword string `mapstructure:"maria_password"`
	DriverName      string `mapstructure:"driver_name"`
	DataSourceName  string `mapstructure:"data_source_name"`
}
