package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

var initialDelta int64
var maxCPUPerReplica int64
var minimumCPUValue int64
var minimumDelta int64

var bnv2Cmd = &cobra.Command{
	Use:   "bnv2",
	Short: "the bnv2 algorithm",
	Long:  "the bnv2 algorithm",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		autoConfigAgent, err := strategies.NewBNV2(initialDelta, GetEndpoints(), GetResources(), initialCPU, initialMemory, maxCPUPerReplica, minimumCPUValue, minimumDelta)
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

	bnv2Cmd.Flags().Int64Var(&initialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	bnv2Cmd.MarkFlagRequired("initcpu")

	bnv2Cmd.Flags().Int64Var(&initialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	bnv2Cmd.MarkFlagRequired("initmem")

	bnv2Cmd.Flags().Int64Var(&maxCPUPerReplica, "maxcpuperreplica", 0, "")
	bnv2Cmd.MarkFlagRequired("maxcpuperreplica")

	bnv2Cmd.Flags().Int64Var(&initialDelta, "initialdelta", 0, "")
	bnv2Cmd.MarkFlagRequired("initialdelta")

	bnv2Cmd.Flags().Int64Var(&minimumCPUValue, "mincpu", 0, "")
	bnv2Cmd.MarkFlagRequired("mincpu")

	bnv2Cmd.Flags().Int64Var(&minimumDelta, "mindelta", 0, "")
	bnv2Cmd.MarkFlagRequired("mindelta")

	AutoConfigCmd.AddCommand(bnv2Cmd)
}
