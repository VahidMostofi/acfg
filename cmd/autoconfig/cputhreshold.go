package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

var indicator string
var threshold float64
var initialCPU int64
var initialMemory int64

var cpuThresholdCmd = &cobra.Command{
	Use:   "cput",
	Short: "increase replica if mean CPU utilization passes threshold",
	Long:  "increase replica if the CPU utilization is more than a specific amount",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		autoConfigAgent, err := strategies.NewCPUThreshold(indicator, threshold, GetEndpoints(), GetResources(), initialCPU, initialMemory)
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

func init() {
	cpuThresholdCmd.Flags().StringVar(&indicator, "indicator", "mean", "indicator of CPU Utilization (default mean).")

	cpuThresholdCmd.Flags().Float64Var(&threshold, "threshold", 0, "threshold value.")
	cpuThresholdCmd.MarkFlagRequired("threshold")

	cpuThresholdCmd.Flags().Int64Var(&initialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	cpuThresholdCmd.MarkFlagRequired("initcpu")

	cpuThresholdCmd.Flags().Int64Var(&initialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	cpuThresholdCmd.MarkFlagRequired("initmem")

	AutoConfigCmd.AddCommand(cpuThresholdCmd)
}
