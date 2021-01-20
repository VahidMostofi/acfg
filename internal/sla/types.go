package sla

import (
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type SLA struct{
	Conditions []Condition
}

type Condition struct{
	Type string					`yaml:"type"`
	EndpointName string			`yaml:"endpointName"`
	Threshold float64			`yaml:"threshold"`
	ComputeFnName string		`yaml:"computeFunctionName"`
}

func (c *Condition) GetComputeFunction() func([]float64) float64{
	if strings.ToLower(c.ComputeFnName) == "mean"{
		return func(values []float64) float64{
			m, err := stats.Mean(values)
			if err != nil{
				panic(err)
			}
			return m
		}
	} else if strings.Contains(strings.ToLower(c.ComputeFnName), "percentile_"){
		return func(values []float64) float64 {
			percent, err := strconv.ParseFloat(strings.Replace(c.ComputeFnName,"percentile_", "",1), 64 )
			if err != nil{
				panic(errors.Wrapf(err,"cant parse %s to get ComputeFn function for condition of SLA. Acceptable example is percentile_90", c.ComputeFnName))
			}
			v, err := stats.Percentile(values, percent)
			return v
		}
	}
	panic(errors.New("unknown computeFnName for SLA condition: " + c.ComputeFnName))
}