package sysstructureagg

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"testing"
)

func TestSystemStructure_GetEndpoints2Resources(t *testing.T) {
	e2r := make(map[string][]string)

	e2r["login"] = []string{"gateway", "auth"}
	e2r["get-book"] = []string{"gateway", "books"}
	e2r["edit-book"] = []string{"gateway", "books"}

	viper.Set(constants.CONFIG_ENDPOINTS_2_RESOURCES, e2r)

	ss, err := NewSystemStructure("predefined")
	if err != nil{
		t.Log(err)
		t.Fail()
		return
	}

	for key, value := range ss.GetEndpoints2Resources(){
		fmt.Println(key, value)
	}
}