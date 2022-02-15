package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

var bnv1DeltaCpu int64
var bnv1DeltaMem int64

var bnv1memCmd = &cobra.Command{
	Use:   "bnv1mem",
	Short: "the bnv1 algorithm with PSI scaling of mem + cpu",
	Long:  "the bnv1 algorithm with PSI scaling of mem + cpu",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		autoConfigAgent, err := strategies.NewBNV1Mem(bnv1DeltaCpu, bnv1DeltaMem, GetEndpoints(), GetResources(), initialCPU, initialMemory)
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

	bnv1memCmd.Flags().Int64Var(&initialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	bnv1memCmd.MarkFlagRequired("initcpu")

	bnv1memCmd.Flags().Int64Var(&initialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	bnv1memCmd.MarkFlagRequired("initmem")

	bnv1memCmd.Flags().Int64Var(&bnv1DeltaCpu, "cpudelta", 0, "")
	bnv1memCmd.MarkFlagRequired("cpudelta")

	bnv1memCmd.Flags().Int64Var(&bnv1DeltaMem, "memdelta", 0, "")
	bnv1memCmd.MarkFlagRequired("memdelta")

	AutoConfigCmd.AddCommand(bnv1memCmd)
}
