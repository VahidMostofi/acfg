package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

var initialCPUDelta int64
var initialMemDelta int64
var maxMemPerReplica int64
var minimumMemValue int64
var minimumCpuDelta int64
var minimumMemDelta int64

var bnv2memCmd = &cobra.Command{
	Use:   "bnv2mem",
	Short: "the bnv2 algorithm with PSI scaling of mem + cpu",
	Long:  "the bnv2 algorithm with PSI scaling of mem + cpu",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		autoConfigAgent, err := strategies.NewBNV2Mem(GetEndpoints(), GetResources(), initialCPU, initialMemory, maxCPUPerReplica, maxMemPerReplica, initialCPUDelta, initialMemDelta, minimumCPUValue, minimumMemValue, minimumCpuDelta, minimumMemDelta)
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

	bnv2memCmd.Flags().Int64Var(&initialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	bnv2memCmd.MarkFlagRequired("initcpu")

	bnv2memCmd.Flags().Int64Var(&initialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	bnv2memCmd.MarkFlagRequired("initmem")

	bnv2memCmd.Flags().Int64Var(&maxCPUPerReplica, "maxcpuperreplica", 0, "")
	bnv2memCmd.MarkFlagRequired("maxcpuperreplica")

	bnv2memCmd.Flags().Int64Var(&maxMemPerReplica, "maxmemperreplica", 0, "")
	bnv2memCmd.MarkFlagRequired("maxmemperreplica")

	bnv2memCmd.Flags().Int64Var(&initialCPUDelta, "initialcpudelta", 0, "")
	bnv2memCmd.MarkFlagRequired("initialcpudelta")

	bnv2memCmd.Flags().Int64Var(&initialMemDelta, "initialmemdelta", 0, "")
	bnv2memCmd.MarkFlagRequired("initialmemdelta")

	bnv2memCmd.Flags().Int64Var(&minimumCPUValue, "mincpu", 0, "")
	bnv2memCmd.MarkFlagRequired("mincpu")

	bnv2memCmd.Flags().Int64Var(&minimumMemValue, "minmem", 0, "")
	bnv2memCmd.MarkFlagRequired("minmem")

	bnv2memCmd.Flags().Int64Var(&minimumCpuDelta, "mincpudelta", 0, "")
	bnv2memCmd.MarkFlagRequired("mincpudelta")

	bnv2memCmd.Flags().Int64Var(&minimumMemDelta, "minmemdelta", 0, "")
	bnv2memCmd.MarkFlagRequired("minmemdelta")

	AutoConfigCmd.AddCommand(bnv2memCmd)
}
