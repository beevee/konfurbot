package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jessevdk/go-flags"

	"github.com/beevee/konfurbot"
	"github.com/beevee/konfurbot/telegram"
	"github.com/beevee/konfurbot/yaml"
)

var logger log.Logger

func main() {
	var opts struct {
		LogFile       string `short:"l" long:"logfile" description:"log file name (writes to stdout if not specified)" env:"KONFURBOT_LOGFILE"`
		ScheduleFile  string `short:"s" long:"schedulefile" default:"schedule.yml" description:"schedule YAML file name" env:"KONFURBOT_SCHEDULEFILE"`
		CheckMode     bool   `short:"c" long:"checkmode" description:"check schedule file for correctness and exit"`
		TelegramToken string `short:"t" long:"telegram-token" description:"@KonfurBot Telegram token" env:"KONFURBOT_TOKEN"`
		Timezone      string `short:"z" long:"timezone" default:"Asia/Yekaterinburg" description:"Local timezone" env:"KONFURBOT_TIMEZONE"`
	}
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(0)
	}

	if opts.LogFile == "" {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	} else {
		logfile, err := os.OpenFile(opts.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open logfile %s: %s", opts.LogFile, err)
			os.Exit(1)
		}
		defer logfile.Close()
		logger = log.NewLogfmtLogger(log.NewSyncWriter(logfile))
	}
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	logger.Log("msg", "starting program", "pid", os.Getpid())

	tz, err := time.LoadLocation(opts.Timezone)
	if err != nil {
		logger.Log("msg", "failed to recognize timezone", "error", err)
		os.Exit(1)
	}

	scheduleYAMLContents, err := ioutil.ReadFile(opts.ScheduleFile)
	if err != nil {
		logger.Log("msg", "failed to open schedule YAML file", "error", err)
		os.Exit(1)
	}
	scheduleStorage := &konfurbot.Schedule{}
	err = yaml.FillScheduleStorage(scheduleStorage, scheduleYAMLContents, tz)
	if err != nil {
		logger.Log("msg", "failed to parse schedule YAML", "error", err)
		os.Exit(1)
	}
	if opts.CheckMode {
		logger.Log("msg", "successfully parsed schedule YAML, and we are in check mode, exiting now")
		os.Exit(0)
	}

	bot := &telegram.Bot{
		ScheduleStorage: scheduleStorage,
		TelegramToken:   opts.TelegramToken,
		Timezone:        tz,
		Logger:          log.NewContext(logger).With("component", "telegram"),
	}

	mustStart(bot)

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	logger.Log("msg", "received signal", "signal", <-signalChannel)

	mustStop(bot)
}

func mustStart(service konfurbot.Service) {
	name := reflect.TypeOf(service)

	logger.Log("msg", "starting service", "name", name)
	if err := service.Start(); err != nil {
		logger.Log("msg", "error starting service", "name", name, "error", err)
		os.Exit(1)
	}
	logger.Log("msg", "started service", "name", name)
}

func mustStop(service konfurbot.Service) {
	name := reflect.TypeOf(service)

	logger.Log("msg", "stopping service", "name", name)
	if err := service.Stop(); err != nil {
		logger.Log("msg", "error stopping service", "name", name, "error", err)
		os.Exit(1)
	}
	logger.Log("msg", "stopped service", "name", name)
}
