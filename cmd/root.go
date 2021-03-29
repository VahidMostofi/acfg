package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/cmd/autoconfig"
	"github.com/vahidmostofi/acfg/cmd/autoscaling"
	dump "github.com/vahidmostofi/acfg/cmd/dumper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "acfg",
	Short: "acfg is a tool to configure different aspects of a microservice application automatically.",
	Long:  `A tool to configure different aspects of a microservice application automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func init() {
	cobra.OnInitialize(initConfigFile)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file, required, would be overwrittern by envs.")
	rootCmd.MarkPersistentFlagRequired("config")

	rootCmd.AddCommand(autoconfig.AutoConfigCmd)
	rootCmd.AddCommand(autoscaling.AutoScaleCmd)
	rootCmd.AddCommand(dump.DumperCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfigFile() {
	if cfgFile == "" {
		fmt.Println("you must pass the config. use --config")
		return
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
		fmt.Println("intializing config file is done.")
	}
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
