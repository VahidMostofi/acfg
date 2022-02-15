package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

// TODO combine CPU & mem threshold
var memIndicator string
var memThreshold float64
var memInitialCPU int64
var memInitialMemory int64

var memThresholdCmd = &cobra.Command{
	Use:   "memt",
	Short: "increase replica if mean mem utilization passes threshold",
	Long:  "increase replica if the mem utilization is more than a specific amount",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		autoConfigAgent, err := strategies.NewMemThreshold(memIndicator, memThreshold, GetEndpoints(), GetResources(), memInitialCPU, memInitialMemory)
		if err != nil {
			panic(err)
		}
		viper.Set(constants.StrategyName, autoConfigAgent.GetName())

		acfgManager, err := factory.NewAutoConfigureManager()
		if err != nil {
			panic(err)
		}
		wl := workload.GetTargetWorkload()
		err = acfgManager.Run(viper.GetString(constants.TestName), autoConfigAgent, &wl)
		if err != nil {
			panic(err)
		}
	},
}

// DO NOT USE WITHOUT FIXING
func init() {
	memThresholdCmd.Flags().StringVar(&memIndicator, "indicator", "mean", "indicator of mem Utilization (default mean).")

	memThresholdCmd.Flags().Float64Var(&memThreshold, "threshold", 0, "threshold value.")
	memThresholdCmd.MarkFlagRequired("threshold")

	memThresholdCmd.Flags().Int64Var(&memInitialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	memThresholdCmd.MarkFlagRequired("initcpu")

	memThresholdCmd.Flags().Int64Var(&memInitialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	memThresholdCmd.MarkFlagRequired("initmem")

	AutoConfigCmd.AddCommand(memThresholdCmd)
}
