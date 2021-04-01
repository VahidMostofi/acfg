package autoconfig

import (
	"github.com/spf13/cobra"
	"github.com/vahidmostofi/acfg/internal/factory"
)

var AutoConfigCmd = &cobra.Command{
	Use:   "autoconfig",
	Short: "autoconfig runs the autoconfiguration",
	Long:  `runs autoconfiguration`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {

}

func getEndpoints() []string {
	t, err := factory.GetEndpointsFilters()
	if err != nil {
		panic(err)
	}
	res := make([]string, 0)
	for s := range t {
		res = append(res, s)
	}
	return res
}

func getResources() []string {
	t, err := factory.GetResourceFilters()
	if err != nil {
		panic(err)
	}
	res := make([]string, 0)
	for s := range t {
		res = append(res, s)
	}
	return res
}
