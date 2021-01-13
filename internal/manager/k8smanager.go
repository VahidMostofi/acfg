package manager

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v12 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
	"sync"
	"time"
)

type Manager interface{}

type K8sManager struct{
	clientSet *kubernetes.Clientset
	namespace string
}

func (k *K8sManager) Deploy(ctx context.Context, reader io.Reader) error{

	b, err := ioutil.ReadAll(reader)
	if err != nil{
		return errors.Wrap(err, "error while reading the reader")
	}

	obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(b, nil, nil)
	if err != nil{
		return errors.Wrap(err, "error while decoding using UniversalSerializer")
	}

	switch o := obj.DeepCopyObject().(type) {
	case *v1.Deployment:
		ptr := obj.DeepCopyObject().(*v1.Deployment)
		_, err = k.clientSet.AppsV1().Deployments(k.namespace).Create(ctx, ptr, metav1.CreateOptions{})
		if err != nil{
			return errors.Wrapf(err, "error while creating Deployment %s", ptr.Name)
		}
	default:
		return errors.Errorf("unknown object: %v", o)
	}
	return nil
}

func (k *K8sManager) UpdateConfigurationsAndWait(ctx context.Context) error{
	// there is a configuration object
	// create the deployment using that, d would be the deployment object
	// for each configuration we update we call wg.Add(1)
	var d *v1.Deployment = nil
	err := k.updateDeployment(ctx, d)
	if err != nil{
		return errors.Wrapf(err, "error while updating deployment for %s", d.Name)
	}
	//go waitDeploymentHaveDesiredCondition(context.TODO(), deploymentsClient,"NewReplicaSetAvailable", targetedDeployments[s].Name,wg)
	//wg.Done()

	return nil //TODO return appropriate error
}

func (k *K8sManager) updateDeployment(ctx context.Context, targetDeployment *v1.Deployment) error {
	log.Debug("updateDeployment() updating deployment", targetDeployment.Name)
	deploymentsClient := k.clientSet.AppsV1().Deployments(k.namespace)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(ctx, targetDeployment.Name, metav1.GetOptions{})
		if getErr != nil{
			errors.Wrap(getErr, fmt.Sprintf("failed to get latest version of Deployment: %s", targetDeployment.Name))
		}
		beforeBytes, _ := result.Marshal()

		returned, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		afterBytes, _ := returned.Marshal()

		isChanged := bytes.Compare(afterBytes,beforeBytes) != 0
		log.Debugf("updateDeployment() comparing before and after for deployment %s: %t", targetDeployment.Name, isChanged)
		if isChanged{
			waitDeploymentHaveDesiredCondition(ctx, deploymentsClient, "ReplicaSetUpdated", targetDeployment.Name,nil,10 * time.Second) //TODO 10 second is hard coded
		}

		return updateErr
	})

	return errors.Wrapf(retryErr, "error while updating deployment: %s", targetDeployment.Name)
}


func waitDeploymentHaveDesiredCondition (ctx context.Context, deploymentClient v12.DeploymentInterface, desiredReason string, deploymentName string, wg *sync.WaitGroup, interval time.Duration){
	wait.Poll(interval, 50*time.Minute, func() (bool, error){
		var isNewReplicaSetAvailable bool
		dep, err := deploymentClient.Get(ctx, deploymentName, metav1.GetOptions{})
		if err != nil{
			panic(err)
		}

		for _,c := range dep.Status.Conditions {
			isNewReplicaSetAvailable = c.Reason == desiredReason
		}

		return isNewReplicaSetAvailable, nil
	})
	if wg != nil{
		wg.Done()
	}
}
func int32Ptr(i int32) *int32 { return &i }
