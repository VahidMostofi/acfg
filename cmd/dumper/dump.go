package dump

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// "github.com/vahidmostofi/acfg/internal/constants"
	"os"

	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/historyextractor"
)

var (
	// cfgFile    string
	outputPath   string
	startTime    int64
	finishTime   int64
	duration     string
	withCpuUtils bool
)

var DumperCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump history info about the system",
	Long:  `dump history info about the system`,
	Run: func(cmd *cobra.Command, args []string) {
		if duration != "" {
			d, err := time.ParseDuration(duration)
			if err != nil {
				panic(err)
			}
			fmt.Println("duration", d.Seconds())
			startTime = finishTime - int64(d.Seconds())
		} else if startTime > -1 {

		} else {
			log.Error("Either duration or the startTime should be passed.")
			os.Exit(1)
		}

		if startTime >= finishTime {
			panic("startTime cant be larger than finishTime.")
		}

		viper.Set(constants.DumpOutputPath, outputPath)
		viper.Set(constants.DumpStartTime, startTime)
		viper.Set(constants.DumpFinishTime, finishTime)
		viper.Set(constants.DumpWithCPUInfo, withCpuUtils)

		historyextractor.DumpHistory()
	},
}

func init() {
	DumperCmd.Flags().StringVar(&outputPath, "output", "", "dumps output path with .json as the format.")
	DumperCmd.MarkFlagRequired("output")

	DumperCmd.Flags().Int64Var(&finishTime, "finishTime", time.Now().Unix(), "finish time of the dump in Unix seconds. Default to now.")
	DumperCmd.Flags().Int64Var(&startTime, "startTime", -1, "start time of the dump in Unix seconds. Either this or the duration should be passed.")
	DumperCmd.Flags().StringVar(&duration, "duration", "", "duration of the queries for dump (24h, 10m, ...). Either this or the startTime should be passed.")

	DumperCmd.Flags().BoolVar(&withCpuUtils, "cpu", false, "extract CPU utilization info too or not (it is slow).")
}
