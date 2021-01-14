package dataaccess

import "github.com/vahidmostofi/acfg/internal/autocfg"

type ConfigDatabase interface{
	Store(code string, data *autocfg.AggregatedData) error
	Retrieve(code string) (*autocfg.AggregatedData, error) // if there is no config with this hash returns nil,false
}
