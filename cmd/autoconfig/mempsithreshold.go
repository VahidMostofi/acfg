package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

var memPSIindicator string
var memPSIthreshold float64
var memPSIinitialCPU int64
var memPSIinitialMemory int64

var memPSIcpuThresholdCmd = &cobra.Command{
	Use:   "mempsit",
	Short: "increase replica if mean CPU PSI utilization passes threshold",
	Long:  "increase replica if the CPU PSI utilization is more than a specific amount",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		// BO
		//autoConfigAgent, err := strategies.NewMemPSIThreshold(memPSIindicator, memPSIthreshold, GetEndpoints(), GetResources(), memPSIinitialCPU, memPSIinitialMemory)
		// Simple
		autoConfigAgent, err := strategies.NewMemPSIThreshold_2(memPSIindicator, memPSIthreshold, GetEndpoints(), GetResources(), memPSIinitialCPU, memPSIinitialMemory)
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
	memPSIcpuThresholdCmd.Flags().StringVar(&memPSIindicator, "indicator", "mean", "indicator of CPU Utilization (default mean).")

	memPSIcpuThresholdCmd.Flags().Float64Var(&memPSIthreshold, "threshold", 0, "threshold value.")
	memPSIcpuThresholdCmd.MarkFlagRequired("threshold")

	memPSIcpuThresholdCmd.Flags().Int64Var(&memPSIinitialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	memPSIcpuThresholdCmd.MarkFlagRequired("initcpu")

	memPSIcpuThresholdCmd.Flags().Int64Var(&memPSIinitialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	memPSIcpuThresholdCmd.MarkFlagRequired("initmem")

	AutoConfigCmd.AddCommand(memPSIcpuThresholdCmd)
}
