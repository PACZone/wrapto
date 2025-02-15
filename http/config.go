package http

type Config struct {
	Port     string `yaml:"port"`
	LockAddr string `yaml:"lock_address"`
}
