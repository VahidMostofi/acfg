package main

import (
	"github.com/vahidmostofi/acfg/cmd"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
    lvl, ok := os.LookupEnv("LOG_LEVEL")
    // LOG_LEVEL not set, let's default to debug
    if !ok {
        lvl = "debug"
    }
    // parse string, this is built-in feature of logrus
    ll, err := logrus.ParseLevel(lvl)
    if err != nil {
        ll = logrus.DebugLevel
    }
    // set global log level
    logrus.SetLevel(ll)
}

func main() {
	cmd.Execute()
	//log.SetReportCaller(true)
}