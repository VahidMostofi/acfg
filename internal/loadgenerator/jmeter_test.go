package loadgenerator

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestJMeterLocalLoadGenerator_Start(t *testing.T) {
	err2 := godotenv.Load("../../test.env")
	if err2 != nil{
		panic(err2)
	}
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("ACFG")
	fmt.Println(viper.Get("ENDPOINTSAGGREGATOR.ARGS.TOKEN"))
	log.SetLevel(log.DebugLevel)
	endpointsFilters := map[string]map[string]interface{}{
		"home":{"uri_regex": "tools.descartes.teastore.webui\\/\"$","http_method": "GET"},
		"login":{"uri_regex": "tools.descartes.teastore.webui\\/loginAction","http_method": "POST"},
		"list-products":{"uri_regex": "tools.descartes.teastore.webui\\/category","http_method": "GET"},
		"look-at-product":{"uri_regex": "tools.descartes.teastore.webui\\/product*","http_method": "GET"},
		"add-to-cart":{"uri_regex": "tools.descartes.teastore.webui\\/cartAction","http_method": "POST"},
	}
	
	wa, err := newWorkloadAggregator()
	if err != nil{
		panic(err)
	}
	fmt.Println(wa)
	// ---------------------
	var k = &JMeterLocalDocker{Command:"-t /input.jmx -Jhostname 172.21.0.2 -Jport 9099 -JnumUser 1 -JrampUp 1 -l mylogfile.log -n"}
	f, err := os.Open("/home/vahid/workspace/t/teastore_browse_nogui_.jmx")
	if err != nil{
		panic(err)
	}
	d, err := ioutil.ReadAll(f)
	k.Data = d
	// ---------------------
	for i := 0;i < 4; i++{
		err = k.Start(nil)
		if err != nil{
			panic(err)
			t.Fail()
			return
		}
		time.Sleep(180 * time.Second)
		// ---------------------
		startTime := time.Now().Unix()
		time.Sleep(180 * time.Second)
		finishTime := time.Now().Unix()
		// ---------------------
		err = k.Stop()
		if err != nil{
			panic(err)
			t.Fail()
			return
		}
		// ---------------------
		time.Sleep(15 * time.Second)
		// ---------------------
		w, err := wa.GetWorkload(startTime, finishTime, endpointsFilters)
		if err != nil{
			panic(err)
		}
		fmt.Println(w.String())
		// ---------------------
		fb,_ := k.GetFeedback()
		if fb != nil{
			fmt.Println(fb["data"])
		}
	}
}
