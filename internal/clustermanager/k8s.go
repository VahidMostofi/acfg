package clustermanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/constants"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v12 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

type K8s struct {
	clientSet        *kubernetes.Clientset
	namespace        string
	deploymentsNames []string
}

func NewK8ClusterManager(namespace string, deploymentNames []string) (ClusterManager, error) {
	kubeconfig := viper.GetString(constants.KubeConfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	k := &K8s{
		namespace:        namespace,
		deploymentsNames: deploymentNames,
		clientSet:        clientset,
	}

	return k, nil
}

func (k *K8s) WaitAllDeploymentsAreStable(ctx context.Context) {
	log.Debug("WaitAllDeploymentsAreStable() waiting for all deployments to be available.")
	deploymentsClient := k.clientSet.AppsV1().Deployments(k.namespace)
	wg := &sync.WaitGroup{}
	for _, deploymentName := range k.deploymentsNames {
		wg.Add(1)
		waitDeploymentHaveDesiredCondition(ctx, deploymentsClient, "NewReplicaSetAvailable", deploymentName, wg, 10)
	}
	wg.Wait()
	log.Debug("WaitAllDeploymentsAreStable() all deployments to be available.")
}

func (k *K8s) Deploy(ctx context.Context, reader io.Reader) error {

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "error while reading the reader")
	}

	obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(b, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error while decoding using UniversalSerializer")
	}

	switch o := obj.DeepCopyObject().(type) {
	case *v1.Deployment:
		ptr := obj.DeepCopyObject().(*v1.Deployment)
		_, err = k.clientSet.AppsV1().Deployments(k.namespace).Create(ctx, ptr, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrapf(err, "error while creating Deployment %s", ptr.Name)
		}
	default:
		return errors.Errorf("unknown object: %v", o)
	}
	return nil
}

func (k *K8s) UpdateConfigurationsAndWait(ctx context.Context, config map[string]*configuration.Configuration) error {

	wg := &sync.WaitGroup{}
	fmt.Println("namespace is", k.namespace)
	deploymentsClient := k.clientSet.AppsV1().Deployments(k.namespace)
	for resourceName, c := range config {
		log.Debugf("UpdateConfigurationsAndWait() updating deployment %s with %s", resourceName, c.String())
		deploymentObj, getErr := deploymentsClient.Get(ctx, resourceName, metav1.GetOptions{})
		if getErr != nil {
			log.Errorf("%s", getErr.Error()) //TODO this is not enoough
			return errors.Wrap(getErr, fmt.Sprintf("failed to get latest version of Deployment: %s", resourceName))
		}

		// update replica count
		deploymentObj.Spec.Replicas = int32Ptr(int32(*c.ReplicaCount))
		// update CPU and memory
		deploymentObj.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				"cpu":    resource.MustParse(c.GetCPUStringForK8s()),
				"memory": resource.MustParse(c.GetMemoryStringForK8s()),
			},
			Requests: corev1.ResourceList{
				"cpu":    resource.MustParse(c.GetCPUStringForK8s()),
				"memory": resource.MustParse(c.GetMemoryStringForK8s()),
			},
		}

		//TODO patch environment variables

		wg.Add(1)
		log.Debugf("UpdateConfigurationsAndWait() calling updateDeployment deployment %s", resourceName)
		err := k.updateDeployment(ctx, deploymentObj)
		if err != nil {
			return errors.Wrapf(err, "error while updating deployment for %s", deploymentObj.Name)
		}
		go waitDeploymentHaveDesiredCondition(ctx, deploymentsClient, "NewReplicaSetAvailable", deploymentObj.Name, wg, time.Second*10)
	}

	wg.Wait()

	return nil
}

func (k *K8s) updateDeployment(ctx context.Context, targetDeployment *v1.Deployment) error {
	log.Debug("updateDeployment() updating deployment", targetDeployment.Name)
	deploymentsClient := k.clientSet.AppsV1().Deployments(k.namespace)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(ctx, targetDeployment.Name, metav1.GetOptions{})
		if getErr != nil {
			return errors.Wrap(getErr, fmt.Sprintf("failed to get latest version of Deployment: %s", targetDeployment.Name))
		}
		beforeBytes, _ := result.Marshal()

		returned, updateErr := deploymentsClient.Update(ctx, targetDeployment, metav1.UpdateOptions{})
		afterBytes, _ := returned.Marshal()

		isChanged := bytes.Compare(afterBytes, beforeBytes) != 0
		log.Debugf("updateDeployment() comparing before and after for deployment %s: %t", targetDeployment.Name, isChanged)
		if isChanged {
			//waitDeploymentHaveDesiredCondition(ctx, deploymentsClient, "ReplicaSetUpdated", targetDeployment.Name,nil,10 * time.Second) //TODO 10 second is hard coded
		}

		return updateErr
	})

	return errors.Wrapf(retryErr, "error while updating deployment: %s", targetDeployment.Name)
}

// waitDeploymentHaveDesiredCondition only works with a single container ! //TODO
func waitDeploymentHaveDesiredCondition(ctx context.Context, deploymentClient v12.DeploymentInterface, desiredReason string, deploymentName string, wg *sync.WaitGroup, interval time.Duration) {
	log.Debugf("wating for container with condition %s for deployment %s", desiredReason, deploymentName)
	wait.Poll(interval, 50*time.Minute, func() (bool, error) {
		var flag bool
		dep, err := deploymentClient.Get(ctx, deploymentName, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		for _, c := range dep.Status.Conditions {
			log.Debugf("for %s got %s", deploymentName, c.Reason)
			flag = c.Reason == desiredReason
			if flag {
				break
			}
		}
		if flag {
			log.Debugf("container with desired condition %s found for deployment %s", desiredReason, deploymentName)
		}

		return flag, nil
	})
	log.Debugf("done wating for container with condition %s for deployment %s", desiredReason, deploymentName)
	if wg != nil {
		wg.Done()
	}
}
func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
