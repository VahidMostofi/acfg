package autocfg

import (
	"fmt"
	sla2 "github.com/vahidmostofi/acfg/internal/sla"
	"gopkg.in/yaml.v2"
	"testing"
)

func Test_One(t *testing.T) {
	sla := &sla2.SLA{
		Conditions: []sla2.Condition{
			{"ResponseTime", "login", 0.250, "percentile_90"},
			{"ResponseTime", "get-book", 0.250, "percentile_90"},
			{"ResponseTime", "edit-book", 0.250, "percentile_90"},
		},
	}
	b, err := yaml.Marshal(sla)
	if err != nil{
		panic(err)
	}
	fmt.Println(string(b))
}
