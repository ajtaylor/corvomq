package main

import (
	"runtime"
	"time"

	"github.com/ajtaylor/corvomq/queue/workers"
	"github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func init() {
	rl, _ := rotatelogs.New("/var/log/corvomq/queue.%Y%m%d%",
		rotatelogs.WithClock(rotatelogs.UTC),
		rotatelogs.WithRotationTime(time.Hour*24))
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(rl)
}

func main() {
	go workers.RunCreateMQVMSubscriber()
	runtime.Goexit()
}
