package logger

type Config struct {
	Filename   string   `yaml:"file_name"`
	LogLevel   string   `yaml:"log_level"`
	Targets    []string `yaml:"targets"`
	MaxSize    int      `yaml:"max_size"`
	MaxBackups int      `yaml:"max_backups"`
	Compress   bool     `yaml:"compress"`
}
