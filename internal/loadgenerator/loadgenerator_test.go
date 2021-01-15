package loadgenerator

import (
	"os"
	"testing"
	"time"
)

func TestK6LocalLoadGenerator_Start(t *testing.T) {
	var k = &K6LocalLoadGenerator{}
	f, err := os.Open("/home/vahid/Desktop/temp.js")
	if err != nil{
		panic(err)
	}

	err = k.Start(nil, f, make(map[string]string))
	if err != nil{
		panic(err)
		t.Fail()
		return
	}
	time.Sleep(140)
}