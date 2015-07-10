package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/ultimatums/ultimate/config"
	"github.com/ultimatums/ultimate/fetch"
	"github.com/ultimatums/ultimate/outputs"
	"github.com/upmio/horus/log"
	"github.com/upmio/horus/utils"
)

const (
	VERSION  = "0.0.1"
	PID_FILE = "ultimate.pid"
	LOG_FILE = "ultimate.log"
	APP_NAME = "ultimate"
)

var (
	cfgFileFlag  = flag.String("config.file", "ultimate.yml", "The configuration file.")
	versionFlag  = flag.Bool("version", false, "Print the version number.")
	reloadFlag   = flag.Bool("reload", false, "Reload the configuration file.")
	logstashFlag = flag.Bool("log.logstash", false, "Generates json in logstash format.")
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *logstashFlag {
		log.ChangeToLogstashFormater(APP_NAME)
	}
	log.SetLogFile(LOG_FILE)

	if *reloadFlag {
		notifyReload()
		os.Exit(0)
	}

	pid := os.Getpid() // This process's pid.
	log.Infof("The process id is %d", pid)
	// Save the pid into the pid file.
	err := utils.WriteFile(PID_FILE, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Fatalln("Error writing this process id:", err)
	}
	defer os.Remove(PID_FILE)

	publisher := outputs.NewPublisherType()
	publisher.Init()

	taskManager := fetch.NewTaskManager(publisher)

	if !reloadConfig(*cfgFileFlag, taskManager) {
		os.Exit(1)
	}

	// Wait for receive a singal to reload configuration file.
	hupCh := make(chan os.Signal)
	hupReady := make(chan bool)
	signal.Notify(hupCh, syscall.SIGHUP)
	go func() {
		<-hupReady
		for range hupCh {
			reloadConfig(*cfgFileFlag, taskManager)
		}
	}()

	go taskManager.Run()
	defer taskManager.Stop()

	close(hupReady)

	// Wait for quit signal to exit this process.
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Kill, syscall.SIGKILL)
	select {
	case s := <-sigCh: // Block until a signal is received.
		log.Warnf("Caught Signal: %v, shuting down gracefully...", s)
	}
	close(hupCh)

	log.Info("See you again!")
	return
}

func notifyReload() {
	pid_b, err := utils.ReadFile(PID_FILE)
	if err != nil {
		log.Fatalln("Error reading the file of this process id:", err)
	}
	pid_s := string(pid_b)

	log.Debug("This process id is ", pid_s)
	if strings.TrimSpace(pid_s) == "" {
		log.Fatalln("The file of this process id is empty.")
	}

	pid, err := strconv.Atoi(pid_s)
	if err != nil {
		log.Fatalf("String to int error: %s", err)
	}

	_, err = utils.ExecCommand(false, "kill", "-HUP", pid_s)
	if err != nil {
		log.Fatalf("Execute commmand <kill -HUP %s> error: %s", pid_s, err)
	}

	err = syscall.Kill(pid, syscall.SYS_READ)
	if err != nil {
		log.Fatalf("Kill signal send failed. Pid: %s, error: %s", pid_s, err)
	} else {
		log.Info("Configuration file is reloading...")
	}
}

func reloadConfig(filename string, taskManager *fetch.TaskManager) bool {
	log.Infof("Loading configuration file %s", filename)

	cfg, err := config.LoadConfig(filename)
	if err != nil {
		log.Errorf("Failed to load configuration file (-config.file=%s): %v", filename, err)
		return false
	}

	taskManager.ApplyConfig(cfg)
	return true
}
