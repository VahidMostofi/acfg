package extractor

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// "github.com/vahidmostofi/acfg/internal/constants"
	"os"
	"path/filepath"
	"strings"

	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/historyextractor"
)

var (
	cfgFile string
)

var ExtractorCmd = &cobra.Command{
	Use:   "extract",
	Short: "extract history info about the system",
	Long:  `extract history info about the system`,
	Run: func(cmd *cobra.Command, args []string) {
		historyextractor.DumpHistory()
	},
}

func init() {
	cobra.OnInitialize(initAutoScalerCmd)
	ExtractorCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (there is not default)")
	ExtractorCmd.MarkPersistentFlagRequired("config")

}

func initAutoScalerCmd() {
	if cfgFile == "" {
		fmt.Println("you must pass the config. use --config")
		os.Exit(1)
	} else {
		cfgFile = checkConfigFile(getAbsPathOfConfigFile(cfgFile))
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err == nil {
			log.Info("Using config file:", viper.ConfigFileUsed())
		} else {
			panic(err)
		}

		viper.SetEnvPrefix("ACFG")
		viper.AutomaticEnv()

		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.MergeConfigMap(viper.AllSettings())
	}
	fmt.Println("initAutoScalerCmd Done")
}

func checkConfigFile(in string) string {
	if _, err := os.Stat(in); os.IsNotExist(err) {
		panic(errors.New(fmt.Sprintf("no config file at: %s", in)))
	}
	return in
}

func getAbsPathOfConfigFile(in string) string {
	in = filepath.Clean(in)
	if filepath.IsAbs(in) {
		return in
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	return filepath.Join(dir, in)
}

func getEndpoints() []string {
	t, err := factory.GetEndpointsFilters()
	if err != nil {
		panic(err)
	}
	res := make([]string, 0)
	for s := range t {
		res = append(res, s)
	}
	return res
}

func getResources() []string {
	t, err := factory.GetResourceFilters()
	if err != nil {
		panic(err)
	}
	res := make([]string, 0)
	for s := range t {
		res = append(res, s)
	}
	return res
}
