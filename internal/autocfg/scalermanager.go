package autocfg

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/autoscalers"
)

type AutoScalerManager struct {
	Replicas          map[string]int
	AutoConfigManager *AutoConfigManager //TODO separate the data aggregator object.
}

func (a *AutoScalerManager) Run(testName string, autoscalerAgent autoscalers.Agent, interval int64) error {
	// TODO log things that happen here to a sqllite or something similar.
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	quit := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		json.Unmarshal([]byte(r.URL.Query().Get("value")), &data)
		name := data["resource"].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		b, err := json.Marshal(struct{ TargetReplicas int }{a.Replicas[name]})
		if err != nil {
			panic(err)
		}
		r.Header.Set("Content-Type", "application/json")
		w.Write([]byte(b))
	})
	go http.ListenAndServe(":3333", r)

	go func() {
		for {
			select {
			case <-ticker.C:
				//
				finishTime := time.Now().Unix()
				startTime := time.Now().Unix() - interval
				aggData, err := a.AutoConfigManager.aggregateData(startTime, finishTime)
				if err != nil {
					panic(err) // TODO
				}

				a.Replicas, err = autoscalerAgent.Evaluate(aggData)
				log.Infof("%v", a.Replicas)
				if err != nil {
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
