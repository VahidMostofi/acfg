package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vahidmostofi/acfg/cmd/autoconfig"
	"github.com/vahidmostofi/acfg/cmd/autoscaling"
	"github.com/vahidmostofi/acfg/cmd/extractor"
)

var rootCmd = &cobra.Command{
	Use:   "acfg",
	Short: "acfg is a tool to configure microservice applications automatically",
	Long:  `A tool to configure microservice applications automatically`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	rootCmd.AddCommand(autoconfig.AutoConfigCmd)
	rootCmd.AddCommand(autoscaling.AutoScaleCmd)
	rootCmd.AddCommand(extractor.ExtractorCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
