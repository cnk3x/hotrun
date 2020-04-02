package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.shu.run/fswatch"
	"go.shu.run/log"
)

func main() {
	log.Config("debug", "HotRun", false, true)
	w := fswatch.Start("./")
	defer w.Stop()

	load := newConfig(w, "hotrun.yml")
	go load.Execute()

	w.Handle(load)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	s := <-c
	log.Infof("收到信号: %s", s)
}
