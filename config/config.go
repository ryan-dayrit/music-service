package config

type Config struct {
	Service struct {
		Network string `yaml:"network"`
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
	} `yaml:"service"`

	Database struct {
		DriverName string `yaml:"driver_name"`
		User       string `yaml:"user"`
		DBName     string `yaml:"db_name"`
		SSLMode    string `yaml:"ssl_mode"`
		Password   string `yaml:"password"`
		Host       string `yaml:"host"`
	} `yaml:"database"`
}
