package autoscalers

import (
	"fmt"
	"testing"
)

// func Test_logScalingDecision(t *testing.T) {
// 	h := &Hybrid{}
// 	err := h.logScalingDecision("test", map[string]int{"auth": 1, "books": 2}, map[string]interface{}{"simple": map[string]int{"vahid": 10, "saeed": 20}})
// 	if err != nil {
// 		panic(err)
// 	}
// }

func Test_ReadPredefinedConfigs(t *testing.T) {
	var filePath string = "/home/vahid/Desktop/projects/research-part2/third-party-raw-data/OnlineShoppingStore-WebServerLogs/splits/train_23_test_24/conditions10_15.json"
	h, err := NewHybridAutoscaler(make([]string, 0), make([]string, 0), 50, filePath)
	if err != nil {
		panic(err)
	}
	hh := h.(*Hybrid)
	// fmt.Println(hh.predefinedReplicas)
	fmt.Println(hh.predefinedReplicas[0].WorkloadRange["editbook"].Low)
}
