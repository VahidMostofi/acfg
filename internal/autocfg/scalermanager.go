package autocfg

import (
	"time"
	"sync"
	"fmt"
	"github.com/vahidmostofi/acfg/internal/autoscalers"
)

func (a *AutoConfigManager) RunAutoscaler(testName string, autoscalerAgent autoscalers.Agent, interval int64) error {
	// TODO log things that happen here to a sqllite or something similar.
	fmt.Println("hey")
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	quit := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(){
		for{
			select{
			case <- ticker.C:
				//
				finishTime := time.Now().Unix()
				startTime := time.Now().Unix() - interval
				aggData, err := a.aggregatedData(startTime, finishTime)
				if err != nil{
					panic(err) // TODO
				}

				_, err = autoscalerAgent.Evaluate(aggData)
				if err != nil{
					panic(err)
				}
			case <-quit:
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	return nil
}