package workload

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
)

type Workload map[string]string

func GetTargetWorkload() Workload {
	w := viper.GetStringMapString(constants.LoadGeneratorArgs)
	return w
}

func (w *Workload) String() string {
	if w == nil {
		return "workload is nil"
	}
	return fmt.Sprintf("total:%f %v", w.GetTotalCount(), w.GetMap())
}

func (w *Workload) GetTotalCount() float64 {
	if _, ok := w.GetMap()["args_vus"]; ok {
		return -1
	}
	var total float64
	for _, value := range w.GetMap() {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}
		total += v
	}
	return total
}

func (w *Workload) GetMap() map[string]string {
	_w := *w
	wm := map[string]string(_w)
	return wm
}

func (w *Workload) GetMapStringInt() map[string]int {
	var err error
	t := w.GetMap()
	res := make(map[string]int)
	for key, value := range t {
		res[key], err = strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
	}
	return res
}

// TODO
func CompareWorkloads(w1 *Workload, w2 *Workload) int {
	return 0
}
