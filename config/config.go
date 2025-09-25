package config

type Config interface {
	Load() error
	GetConfigPath(env string) (string, error)
	IsSet(key string) bool
	Get(key string) Value
}
