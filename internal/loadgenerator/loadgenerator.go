package loadgenerator

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/workload"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

	removeCmd := exec.Command("docker", "container", "remove", containerName)
	err  = removeCmd.Run()

	if err != nil{
		log.Warnf("error while removing k6 load generator container before running it %v", err)
	}
	fmt.Println(removeCmd.CombinedOutput())

	cmd := exec.Command("docker", "run", "--network", "host", "--rm", "--name", containerName, "-v", file.Name() + ":" + "/script.js", "loadimpact/k6",  "run", "/script.js")
	fmt.Println(cmd.String())
	err = cmd.Run()
	if err != nil{
		out, _ := cmd.CombinedOutput()
		log.Debugf("Failed %s",string(out))
		return errors.Wrapf(err, "error while starting load generator %s", string(out))
	}
	out, _ := cmd.CombinedOutput()
	log.Debugf("Failed %s",string(out))


	return nil
}

func (k *K6LocalLoadGenerator) Stop() error {
	return nil
}

func (k *K6LocalLoadGenerator) GetFeedback() (map[string]interface{}, error) {
	return nil, nil
}

func prepareLoadGenerator(workload *workload.Workload, info map[string]interface{}) ([]byte, error){

	return nil, nil
}
