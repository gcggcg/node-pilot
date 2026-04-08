package config

type Config struct {
	DBPath string
	Listen string
	Debug  bool
}

const (
	OutPutLimitLines = 100
)
