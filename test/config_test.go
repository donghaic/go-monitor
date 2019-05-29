package test

import (
	"gt-monitor/config"
	"testing"
)

func TestConfig(t *testing.T) {
	cnf, _ := config.ReadConfig("../dev-config.yaml")
	println("data file ==", cnf.TaskQueueDataDir)
}
