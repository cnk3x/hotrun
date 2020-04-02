package main

import (
	"regexp"
	"time"

	"go.shu.run/cmd"
	"go.shu.run/fswatch"
)

type RunnerConfig struct {
	Name  string        `yaml:"name"`
	Delay time.Duration `yaml:"delay"`
	Match []string      `yaml:"match"`
	Run   []string      `yaml:"run"`
}

func getRunnerFromConfig(runnerConfig RunnerConfig) (*CommandRunner, error) {
	h := &CommandRunner{}
	if err := h.Update(runnerConfig); err != nil {
		return nil, err
	}
	return h, nil
}

type CommandRunner struct {
	name    string
	delay   time.Duration
	c       *cmd.C
	matches []*regexp.Regexp
}

func (runner *CommandRunner) Delay() time.Duration {
	return runner.delay
}

func (runner *CommandRunner) Update(runnerConfig RunnerConfig) error {
	runner.name = runnerConfig.Name
	if runner.c == nil {
		runner.c = cmd.New(runnerConfig.Run...)
	} else {
		runner.c.Command(runnerConfig.Run...)
	}
	runner.delay = runnerConfig.Delay
	runner.matches = make([]*regexp.Regexp, 0, len(runnerConfig.Match))
	for _, match := range runnerConfig.Match {
		r, err := regexp.Compile(match)
		if err != nil {
			return err
		}
		runner.matches = append(runner.matches, r)
	}
	return nil
}

func (runner *CommandRunner) Name() string {
	return runner.name
}

func (runner *CommandRunner) Match(e fswatch.Event) bool {
	for _, match := range runner.matches {
		if match.MatchString(e.Name) {
			return true
		}
	}
	return false
}

func (runner *CommandRunner) Run() error {
	runner.c.Run()
	return nil
}

func (runner *CommandRunner) Stop() error {
	if runner.c != nil {
		return runner.c.Kill()
	}
	return nil
}

var _ fswatch.Handler = (*CommandRunner)(nil)
