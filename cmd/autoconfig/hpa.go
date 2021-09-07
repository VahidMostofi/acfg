package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

var hpaCmd = &cobra.Command{
	Use:   "hpa",
	Short: "the hpa algorithm",
	Long:  "the hpa algorithm",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		autoConfigAgent, err := strategies.NewHPA(getResources(), threshold, initialCPU, initialMemory)
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

	hpaCmd.Flags().Int64Var(&initialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	hpaCmd.MarkFlagRequired("initcpu")

	hpaCmd.Flags().Int64Var(&initialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	hpaCmd.MarkFlagRequired("initmem")

	hpaCmd.Flags().Float64Var(&threshold, "threshold", 0, "threshold value.")
	hpaCmd.MarkFlagRequired("threshold")

	AutoConfigCmd.AddCommand(hpaCmd)
}
