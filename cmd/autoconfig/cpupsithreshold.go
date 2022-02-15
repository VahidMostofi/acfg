package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

var cpuPSIindicator string
var cpuPSIthreshold float64
var cpuPSIinitialCPU int64
var cpuPSIinitialMemory int64

var cpuPSIcpuThresholdCmd = &cobra.Command{
	Use:   "cpupsit",
	Short: "increase replica if mean CPU PSI utilization passes threshold",
	Long:  "increase replica if the CPU PSI utilization is more than a specific amount",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		autoConfigAgent, err := strategies.NewCPUPSIThreshold(indicator, threshold, GetEndpoints(), GetResources(), initialCPU, initialMemory)
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
	cpuPSIcpuThresholdCmd.Flags().StringVar(&cpuPSIindicator, "indicator", "mean", "indicator of CPU Utilization (default mean).")

	cpuPSIcpuThresholdCmd.Flags().Float64Var(&cpuPSIthreshold, "threshold", 0, "threshold value.")
	cpuPSIcpuThresholdCmd.MarkFlagRequired("threshold")

	cpuPSIcpuThresholdCmd.Flags().Int64Var(&cpuPSIinitialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	cpuPSIcpuThresholdCmd.MarkFlagRequired("initcpu")

	cpuPSIcpuThresholdCmd.Flags().Int64Var(&cpuPSIinitialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	cpuPSIcpuThresholdCmd.MarkFlagRequired("initmem")

	AutoConfigCmd.AddCommand(cpuPSIcpuThresholdCmd)
}
