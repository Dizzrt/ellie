package config

type Config interface {
	Load() error
	Get(key string) Value
	IsSet(key string) bool
	Unmarshal(obj any) error
	UnmarshalKey(key string, obj any) error
	GetConfigPath(env string) (string, error)
}
