package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"go.shu.run/fswatch"
	"go.shu.run/log"

	"gopkg.in/yaml.v2"
)

func newConfig(w *fswatch.Watcher, fn string) *fswatch.Runner {
	return fswatch.NewFunc(&Config{w: w, path: fn})
}

type Config struct {
	w    *fswatch.Watcher
	path string
}

func (ch *Config) Delay() time.Duration {
	return time.Second
}

func (ch *Config) Name() string {
	return "加载配置"
}

func (ch *Config) Match(e fswatch.Event) bool {
	return strings.HasSuffix(e.Name, ch.path)
}

func (ch *Config) Run() error {
	log.Infof("加载配置")
	var cfg C

	v, err := ioutil.ReadFile(ch.path)
	if err != nil {
		log.Errorf("打开配置文件失败: %v", err)
		cfg = ch.writeDefault()
	} else if err = yaml.Unmarshal(v, &cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	log.Config(cfg.Logger.Level, cfg.Logger.Name, cfg.Logger.Caller,cfg.Logger.Time)

	for _, hc := range cfg.HandlerConfigs {
		find := ch.w.Find(func(hf *fswatch.Runner) bool {
			if runner, ok := hf.Handler.(*CommandRunner); ok {
				if runner.Name() == hc.Name {
					if err := runner.Update(hc); err != nil {
						log.Errorf("%v", err)
					}
					return true
				}
			}
			return false
		})

		if find {
			continue
		}

		runner, err := getRunnerFromConfig(hc)
		if err != nil {
			log.Errorf("%v", err)
			continue
		}
		ch.w.Handle(runner)
	}

	return nil
}

func (ch *Config) Stop() error {
	return nil
}

func (ch *Config) writeDefault() C {
	cfg := C{}
	cfg.Logger.Name = "HotRun"
	cfg.Logger.Level = "debug"
	cfg.Logger.Caller = false
	cfg.Logger.Time = true
	v, _ := yaml.Marshal(&cfg)
	_ = ioutil.WriteFile(ch.path, v, 0644)
	return cfg
}

var _ fswatch.Handler = (*Config)(nil)
