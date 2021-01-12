package workload

type Workload map[string]int64

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
