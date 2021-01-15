package loadgenerator

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/workload"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const containerName = "kkkk6localautoconfig"

type LoadGenerator interface{
	Start(workload *workload.Workload, reader io.Reader, extras map[string]string) error
	Stop() error
	GetFeedback() (map[string]interface{}, error)
}

type K6LocalLoadGenerator struct {

}

func (k *K6LocalLoadGenerator) Start(workload *workload.Workload, reader io.Reader, extra map[string]string) error {
	log.Debug("K6LocalLoadGenerator: reading content of script for load generator")
	t, err := ioutil.ReadAll(reader)
	scriptContent := string(t)
	if err != nil{
		return errors.Wrap(err, "error reading script content")
	}

	if extra != nil{
		log.Debugf("K6LocalLoadGenerator: replacing extra information on script content")
		for key,value := range extra{
			log.Debugf("K6LocalLoadGenerator: replacing %s with %s", key, value)
			scriptContent = strings.ReplaceAll(scriptContent, key, value)
		}
	}

	file, err := ioutil.TempFile("", "k6-script")
	if err != nil{
		return errors.Wrap(err, "error while creating temp file for script of k6 local load generator")
	}
	file.WriteString(scriptContent)
	file.Close()
	err = os.Chmod(file.Name(), 0777)
	if err != nil{
		return errors.Wrapf(err, "error while chaning permission of temp file %s", file.Name())
	}

	removeCmd := exec.Command("docker", "container", "rm", "-f", containerName)
	fmt.Println(removeCmd.String())
	err  = removeCmd.Run()

	if err != nil{
		log.Warnf("K6LocalLoadGenerator: error while removing k6 load generator container before running it %v", err)
	}
	b, _ := removeCmd.CombinedOutput()

	log.Warn("K6LocalLoadGenerator: " + string(b))

	cmd := exec.Command("docker", "run", "--network", "host","-d", "--rm", "--name", containerName, "-v", file.Name() + ":" + "/script.js", "loadimpact/k6",  "run", "/script.js")
	log.Warn("K6LocalLoadGenerator: ",cmd.String())
	err = cmd.Run()
	if err != nil{
		out, _ := cmd.CombinedOutput()
		log.Debugf("K6LocalLoadGenerator: Failed %s",string(out))
		return errors.Wrapf(err, "error while starting load generator %s", string(out))
	}
	out, _ := cmd.CombinedOutput()
	log.Debugf("K6LocalLoadGenerator: Failed %s",string(out))
	for{
		url := "http://localhost:6565/v1/status"
		resp, err := http.Get(url)
		if err != nil{
			continue
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil{
			panic(err)
		}
		if strings.Contains(string(b), "\"running\":true"){
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (k *K6LocalLoadGenerator) Stop() error {
	log.Debug("K6LocalLoadGenerator: Stopping load generator.")
	removeCmd := exec.Command("docker", "container", "rm", "-f", containerName)
	err := removeCmd.Run()

	if err != nil{
		log.Warnf("K6LocalLoadGenerator: error while removing k6 load generator container for stopping %v", err)
	}
	b, _ := removeCmd.CombinedOutput()

	log.Warn("K6LocalLoadGenerator: " + string(b))
	return nil
}

func (k *K6LocalLoadGenerator) GetFeedback() (map[string]interface{}, error) {
	return nil, nil
}

func prepareLoadGenerator(workload *workload.Workload, info map[string]interface{}) ([]byte, error){

	return nil, nil
}
