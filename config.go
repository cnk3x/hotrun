package main

type C struct {
	Logger struct {
		Name   string `yaml:"name"`
		Level  string `yaml:"level"`
		Caller bool   `yaml:"caller"`
		Time   bool   `yaml:"time"`
	} `yaml:"log"`
	HandlerConfigs []RunnerConfig `yaml:"handlers"`
}
