package loadgenerator

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/workload"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const jMeterContainerName = "jMeterContainerNameWIthRahadf09283rhfsadf"

type JMeterLocalDocker struct {
	Data []byte
	Args map[string]string
	Command string
	feedback map[string]interface{}
}

func (j* JMeterLocalDocker) Start(workload *workload.Workload, reader io.Reader, extras map[string]string) error{
	log.SetLevel(log.DebugLevel)
	log.Debug("JMeterLocalDocker: reading content of script for load generator")

	scriptContent := string(j.Data)

	if j.Args != nil{
		log.Debugf("JMeterLocalDocker: replacing k.Args information on script content")
		for key,value := range j.Args{
			log.Debugf("JMeterLocalDocker: replacing %s with %s", key, value)
			scriptContent = strings.ReplaceAll(scriptContent, key, value)
		}
	}

	file, err := ioutil.TempFile("", "jmeter-jmx")
	if err != nil{
		return errors.Wrap(err, "error while creating temp file for jmx of jmeter local load generator")
	}
	file.WriteString(scriptContent)
	file.Close()
	err = os.Chmod(file.Name(), 0777)
	if err != nil{
		return errors.Wrapf(err, "jmeter local load generator error while chaning permission of temp file %s", file.Name())
	}

	removeCmd := exec.Command("docker", "container", "rm", "-f", jMeterContainerName)
	err  = removeCmd.Run()

	if err != nil{
		log.Warnf("JMeterLocalDocker: error while removing jmeter load generator container before running it %v", err)
	}
	b, _ := removeCmd.CombinedOutput()

	log.Warnf("JMeterLocalDocker: removing combined output: " + string(b))

	commandParts := strings.Split(j.Command, " ")
	commandArgs := []string{"run", "--network", "host", "-d", "--rm", "--name", containerName, "-v", file.Name() + ":" + "/input.jmx", "vmarrazzo/jmeter"}
	commandArgs = append(commandArgs, commandParts...)
	cmd := exec.Command("docker",  commandArgs...)
	log.Debug("JMeterLocalDocker: ",cmd.String())
	err = cmd.Run()
	if err != nil{
		out, _ := cmd.CombinedOutput()
		log.Debugf("JMeterLocalDocker: failure combined output: %s",string(out))
		return errors.Wrapf(err, "error while starting load generator %s", string(out))
	}
	out, _ := cmd.CombinedOutput()
	log.Debugf("JMeterLocalDocker: running combined output: %s",string(out))
	return nil
}

func (j* JMeterLocalDocker) Stop() error{
	log.Debug("JMeterLocalDocker: Stopping load generator.")
	removeCmd := exec.Command("docker", "container", "rm", "-f", containerName)
	err := removeCmd.Run()

	if err != nil{
		log.Warnf("JMeterLocalDocker: error while removing jmeter load generator container for stopping %v", err)
	}
	b, _ := removeCmd.CombinedOutput()

	log.Warn("JMeterLocalDocker: " + string(b))
	return nil
}

func (j* JMeterLocalDocker) GetFeedback() (map[string]interface{}, error){
	log.Warnf("currently there is no way to get feedback from JMeter")
	return nil,nil
}

