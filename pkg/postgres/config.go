package postgres

type Config struct {
	DriverName string `yaml:"driver_name"`
	User       string `yaml:"user"`
	DBName     string `yaml:"db_name"`
	SSLMode    string `yaml:"ssl_mode"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
}
