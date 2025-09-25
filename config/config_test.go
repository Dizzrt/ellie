package config

import (
	"fmt"
	"testing"
)

func TestStdViperConfig(t *testing.T) {
	c := NewStdViperConfig(
		EnvPrefix("ELLIE"),
	)

	if err := c.Load(); err != nil {
		t.Fatalf("load config failed, err: %v", err)
	}

	port := c.Get("clusters.a.port").Int()
	fmt.Printf("port: %d\n", port)
}
