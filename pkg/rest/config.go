package rest

type Config struct {
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	ServerUrl    string `yaml:"server_url"`
}
