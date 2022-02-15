package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

var bnv1Delta int64

var bnv1Cmd = &cobra.Command{
	Use:   "bnv1",
	Short: "the bnv1 algorithm",
	Long:  "the bnv1 algorithm",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		autoConfigAgent, err := strategies.NewBNV1(bnv1Delta, GetEndpoints(), GetResources(), initialCPU, initialMemory)
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

	bnv1Cmd.Flags().Int64Var(&initialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	bnv1Cmd.MarkFlagRequired("initcpu")

	bnv1Cmd.Flags().Int64Var(&initialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	bnv1Cmd.MarkFlagRequired("initmem")

	bnv1Cmd.Flags().Int64Var(&bnv1Delta, "delta", 0, "")
	bnv1Cmd.MarkFlagRequired("delta")

	AutoConfigCmd.AddCommand(bnv1Cmd)
}
