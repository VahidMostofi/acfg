package loadgenerator

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/workload"
	"k8s.io/apimachinery/pkg/util/json"
)

type K6LocalLoadGenerator struct {
	Data     []byte
	Args     map[string]string
	feedback map[string]interface{}
}

func (k *K6LocalLoadGenerator) Start(workload *workload.Workload) error {
	log.Debug("K6LocalLoadGenerator: reading content of script for load generator")

	scriptContent := string(k.Data)

	if k.Args != nil {
		log.Debugf("K6LocalLoadGenerator: replacing k.Args information on script content")
		for key, value := range k.Args {
			log.Debugf("K6LocalLoadGenerator: replacing %s with %s", key, value)
			scriptContent = strings.ReplaceAll(scriptContent, key, value)
		}
	}

	file, err := ioutil.TempFile("", "k6-script")
	if err != nil {
		return errors.Wrap(err, "error while creating temp file for script of k6 local load generator")
	}
	file.WriteString(scriptContent)
	file.Close()
	err = os.Chmod(file.Name(), 0777)
	if err != nil {
		return errors.Wrapf(err, "error while chaning permission of temp file %s", file.Name())
	}

	removeCmd := exec.Command("docker", "container", "rm", "-f", containerName)
	err = removeCmd.Run()

	if err != nil {
		log.Warnf("K6LocalLoadGenerator: error while removing k6 load generator container before running it %v", err)
	}
	b, _ := removeCmd.CombinedOutput()

	log.Warnf("K6LocalLoadGenerator: removing combined output: " + string(b))

	cmd := exec.Command("docker", "run", "--network", "host", "-d", "--rm", "--name", containerName, "-v", file.Name()+":"+"/script.js", "loadimpact/k6", "run", "/script.js")
	log.Debug("K6LocalLoadGenerator: ", cmd.String())
	err = cmd.Run()
	if err != nil {
		out, _ := cmd.CombinedOutput()
		log.Debugf("K6LocalLoadGenerator: failure combined output: %s", string(out))
		return errors.Wrapf(err, "error while starting load generator %s", string(out))
	}
	out, _ := cmd.CombinedOutput()
	log.Debugf("K6LocalLoadGenerator: running combined output: %s", string(out))
	log.Debugf("K6LocalLoadGenerator: waiting for the load generator to start.")
	for {
		time.Sleep(1 * time.Second)
		url := "http://localhost:6565/v1/status"
		resp, err := http.Get(url)
		if err != nil {
			log.Debugf("K6LocalLoadGenerator: error while getting response and waiting for load generator to start %s", err.Error())
			continue
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		if strings.Contains(string(b), "\"running\":true") {
			log.Debugf("K6LocalLoadGenerator: based on http request the load generator is ready.")
			break
		}
	}

	return nil
}

func (k *K6LocalLoadGenerator) retrieveFeedbackBeforeStop() error {
	resp, err := http.Get("http://localhost:6565/v1/metrics")
	if err != nil {
		return errors.Wrapf(err, "error while getting feedback of k6 loadgenerator (local)")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "error while reading response body.")
	}
	json.Unmarshal(body, &k.feedback)
	return nil
}

func (k *K6LocalLoadGenerator) Stop() error {
	log.Debug("K6LocalLoadGenerator: Getting latest metrics.")
	err := k.retrieveFeedbackBeforeStop()
	if err != nil {
		return err
	}
	log.Debug("K6LocalLoadGenerator: Got metrics.")
	log.Debug("K6LocalLoadGenerator: Stopping load generator.")
	removeCmd := exec.Command("docker", "container", "rm", "-f", containerName)
	err = removeCmd.Run()

	if err != nil {
		log.Warnf("K6LocalLoadGenerator: error while removing k6 load generator container for stopping %v", err)
	}
	b, _ := removeCmd.CombinedOutput()

	log.Warn("K6LocalLoadGenerator: " + string(b))
	return nil
}

func (k *K6LocalLoadGenerator) GetFeedback() (map[string]interface{}, error) {
	return k.feedback, nil
}
