package autoconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
)

var (
	cfgFile  string
	testName string
)

var AutoConfigCmd = &cobra.Command{
	Use:   "autoconfig",
	Short: "autoconfig runs the autoconfiguration",
	Long:  `runs autoconfiguration`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	cobra.OnInitialize(initConfigAutoConfigCmd)
	AutoConfigCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (there is not default)")
	AutoConfigCmd.MarkPersistentFlagRequired("config")

	AutoConfigCmd.PersistentFlags().StringVar(&testName, "name", "", "name of the test (there is no default)")
	AutoConfigCmd.MarkPersistentFlagRequired("name")
}

func initConfigAutoConfigCmd() {
	// if cfgFile == ""{
	// 	fmt.Println("you must pass the config. use --config")
	// 	os.Exit(1)
	// } else {
	if cfgFile == "" {
		return
	}
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

	viper.Set(constants.TestName, testName)
	// }
	fmt.Println("initConfigAutoConfigCmd Done")
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
