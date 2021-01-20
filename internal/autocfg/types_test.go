package autocfg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

func Test_One(t *testing.T) {
	sla := &SLA{
		Conditions: []Condition{
			{"ResponseType", "login", 0.250, "percentile_90"},
			{"ResponseType", "get-book", 0.250, "percentile_90"},
			{"ResponseType", "edit-book", 0.250, "percentile_90"},
		},
	}
	b, err := yaml.Marshal(sla)
	if err != nil{
		panic(err)
	}
	fmt.Println(string(b))
}
