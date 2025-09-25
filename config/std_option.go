package config

import "github.com/spf13/viper"

type Option func(*stdViperConfig)

func ConfigFileDir(dir string) Option {
	return func(c *stdViperConfig) {
		c.configFileDir = dir
	}
}

func ConfigFileName(name string) Option {
	return func(c *stdViperConfig) {
		c.configFileName = name
	}
}

func ConfigFileSuffix(suffix string) Option {
	return func(c *stdViperConfig) {
		c.configFileSuffix = suffix
	}
}

func ValidEnvs(envs []string) Option {
	return func(c *stdViperConfig) {
		c.validEnvs = envs
	}
}

func ViperOption(opts ...viper.Option) Option {
	return func(c *stdViperConfig) {
		c.viperOpts = append(c.viperOpts, opts...)
	}
}

func EnvPrefix(prefix string) Option {
	return func(c *stdViperConfig) {
		c.envPrefix = prefix
	}
}

func AutomaticEnv(enabled bool) Option {
	return func(c *stdViperConfig) {
		c.isAutomaticEnvEnabled = enabled
	}
}
