package autoscaling

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/autoscalers"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
)

var allocationsFile string
var numberOfConfigsToUse int
var hpaCpuUPercentageThreshold int64
var intervalSeconds int64

var hybridAutoscalerCmd = &cobra.Command{
	Use:   "hybrid",
	Short: "hybrid combines hpa with pre-configured.",
	Long:  "hybrid combines hpa with pre-configured.",
	Run: func(cmd *cobra.Command, args []string) {

		if numberOfConfigsToUse > 0 && len(allocationsFile) == 0 {
			panic("if usecount > 0, the path to configs file (allocationsFile) should be passed.")
		}

		autoscalingAgent, err := autoscalers.NewHybridAutoscaler(getEndpoints(), getResources(), hpaCpuUPercentageThreshold, allocationsFile, numberOfConfigsToUse)
		if err != nil {
			panic(err)
		}
		viper.Set(constants.AutoScalingApproachName, autoscalingAgent.GetName())

		autoScalerManager, err := factory.NewAutoScalerManager()
		if err != nil {
			panic(err)
		}

		err = autoScalerManager.Run(viper.GetString(constants.TestName), autoscalingAgent, intervalSeconds)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	hybridAutoscalerCmd.Flags().StringVar(&allocationsFile, "allocationsFile", "", "the path to resource allocations for the predefined configurations.")

	hybridAutoscalerCmd.Flags().IntVar(&numberOfConfigsToUse, "usecount", 0, "how many of the configs in the resource allocation file should be used. -1 for all")
	hybridAutoscalerCmd.MarkFlagRequired("usecount")

	hybridAutoscalerCmd.Flags().Int64Var(&hpaCpuUPercentageThreshold, "hpat", 50, "what is the desired CPU utilization when using HPA.")
	hybridAutoscalerCmd.MarkFlagRequired("hpat")

	hybridAutoscalerCmd.Flags().Int64Var(&intervalSeconds, "interval", 10, "how often should check to update the configuration. How often the autoscaler should Evaluate")
	hybridAutoscalerCmd.MarkFlagRequired("interval")

	AutoScaleCmd.AddCommand(hybridAutoscalerCmd)
}
