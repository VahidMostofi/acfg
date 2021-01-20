package workload

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
)

type Workload map[string]int64

func GetTargetWorkload() Workload{
	temp := viper.GetStringMap(constants.TargetSystemWorkloadBody)
	w := make(map[string]int64)
	for k,v := range temp{
		w[k] = int64(v.(int))
	}
	return w
}

func (w *Workload) String() string{
	if w == nil{
		return "workload is nil"
	}
	return fmt.Sprintf("total:%d %v", w.GetTotalCount(), w.GetMap())
}


func (w *Workload) GetTotalCount() int64{
	var total int64
	for _,v := range w.GetMap(){
		total += v
	}
	return total
}

func (w *Workload) GetMap() map[string]int64{
	_w := *w
	wm := map[string]int64(_w)
	return wm
}

// TODO
func CompareWorkloads(w1 *Workload, w2 *Workload) float64{
	return 0
}
