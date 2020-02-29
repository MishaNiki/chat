package templates

// Config ...
type Config struct {
	Root string `json:"templates-root"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}
