package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vahidmostofi/acfg/cmd/autoconfig"
	"github.com/vahidmostofi/acfg/cmd/autoscaling"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "acfg",
	Short: "acfg is a tool to configure microservice applications automatically",
	Long: `A tool to configure microservice applications automatically`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init()  {
	rootCmd.AddCommand(autoconfig.AutoConfigCmd)
	rootCmd.AddCommand(autoscaling.AutoScaleCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}