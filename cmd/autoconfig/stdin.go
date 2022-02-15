package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/factory"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
)

var pythonPath string
var scriptPath string

var stdinCmd = &cobra.Command{
	Use:   "stdin",
	Short: "gets values from stdin of another python code.",
	Long:  "gets values from stdin of another python code.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here

		autoConfigAgent, err := strategies.NewPythonRunner(pythonPath, scriptPath, GetResources(), initialCPU, initialMemory)
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
	stdinCmd.Flags().Int64Var(&initialCPU, "initcpu", 0, "initial CPU to allocate to each replica. Use 1000 for 1 CPU unit.")
	stdinCmd.MarkFlagRequired("initcpu")

	stdinCmd.Flags().Int64Var(&initialMemory, "initmem", 0, "initial memory to allocate to each replica in MB.")
	stdinCmd.MarkFlagRequired("initmem")

	stdinCmd.Flags().StringVar(&pythonPath, "pythonpath", "", "path to python interpreter")
	stdinCmd.MarkFlagRequired("pythonpath")

	stdinCmd.Flags().StringVar(&scriptPath, "scriptpath", "", "path to python interpreter")
	stdinCmd.MarkFlagRequired("scriptpath")

	AutoConfigCmd.AddCommand(stdinCmd)
}
