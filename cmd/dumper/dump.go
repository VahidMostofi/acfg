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
	outputPath string
	startTime  int64
	finishTime int64
	duration   string
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

		historyextractor.DumpHistory()
	},
}

func init() {
	// cobra.OnInitialize(initDumperCmd)
	// DumperCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (there is not default)")
	// DumperCmd.MarkPersistentFlagRequired("config")

	DumperCmd.Flags().StringVar(&outputPath, "output", "", "dumps output path with .json as the format.")
	DumperCmd.MarkFlagRequired("output")

	DumperCmd.Flags().Int64Var(&finishTime, "finishTime", time.Now().Unix(), "finish time of the dump in Unix seconds. Default to now.")
	DumperCmd.Flags().Int64Var(&startTime, "startTime", -1, "start time of the dump in Unix seconds. Either this or the duration should be passed.")
	DumperCmd.Flags().StringVar(&duration, "duration", "", "duration of the queries for dump (24h, 10m, ...). Either this or the startTime should be passed.")
}

// func initDumperCmd() {
// 	if cfgFile == "" {
// 		fmt.Println("you must pass the config. use --config")
// 		return
// 	}
// 	cfgFile = checkConfigFile(getAbsPathOfConfigFile(cfgFile))
// 	viper.SetConfigFile(cfgFile)

// 	if err := viper.ReadInConfig(); err == nil {
// 		log.Info("Using config file:", viper.ConfigFileUsed())
// 	} else {
// 		panic(err)
// 	}

// 	viper.SetEnvPrefix("ACFG")
// 	viper.AutomaticEnv()

// 	replacer := strings.NewReplacer(".", "_")
// 	viper.SetEnvKeyReplacer(replacer)
// 	viper.MergeConfigMap(viper.AllSettings())
// 	fmt.Println("DumperCmd Done")
// }

// func checkConfigFile(in string) string {
// 	if _, err := os.Stat(in); os.IsNotExist(err) {
// 		panic(errors.New(fmt.Sprintf("no config file at: %s", in)))
// 	}
// 	return in
// }

// func getAbsPathOfConfigFile(in string) string {
// 	in = filepath.Clean(in)
// 	if filepath.IsAbs(in) {
// 		return in
// 	}

// 	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
// 	if err != nil {
// 		panic(err)
// 	}

// 	return filepath.Join(dir, in)
// }

// func getEndpoints() []string {
// 	t, err := factory.GetEndpointsFilters()
// 	if err != nil {
// 		panic(err)
// 	}
// 	res := make([]string, 0)
// 	for s := range t {
// 		res = append(res, s)
// 	}
// 	return res
// }

// func getResources() []string {
// 	t, err := factory.GetResourceFilters()
// 	if err != nil {
// 		panic(err)
// 	}
// 	res := make([]string, 0)
// 	for s := range t {
// 		res = append(res, s)
// 	}
// 	return res
// }
