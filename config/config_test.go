package config

import (
	"fmt"
	"testing"
)

type Cluster struct {
	IP   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}

type Conf struct {
	Clusters map[string]Cluster `mapstructure:"clusters"`
}

func TestStdViperConfig(t *testing.T) {
	c := NewStdViperConfig(
		EnvPrefix("ELLIE"),
	)

	if err := c.Load(); err != nil {
		t.Fatalf("load config failed, err: %v", err)
	}

	port := c.Get("clusters.a.port").Int()
	fmt.Printf("port: %d\n", port)

	var conf Conf
	if err := c.UnmarshalKey("clusters", &conf.Clusters); err != nil {
		t.Fatalf("unmarshal config failed, err: %v", err)
	}
	fmt.Printf("conf: %+v\n", conf)
}
