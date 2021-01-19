package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/cmd"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
	"os"
)

func main() {
	cmd.Execute()
	// TODO this needs to be moved where we actually create CPU utilization
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	acfgManager, err := factory.NewAutoConfigureManager()
	if err != nil{
		log.Panic(err.Error())
		os.Exit(1)
	}
	autoConfigAgent,err := strategies.NewCPUThreshold("mean", 50, []string{"login","get-book","edit-book"}, []string{"auth","books","gateway"},500, 512)
	if err != nil{
		panic(err)
	}
	wl := workload.GetWorkload()
	err = acfgManager.Run(viper.GetString(constants.TestName), autoConfigAgent, &wl)
	if err != nil{
		panic(err)
	}
}