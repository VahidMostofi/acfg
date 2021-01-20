package clustermanager

import (
	"context"
	"flag"
	"fmt"
	"github.com/vahidmostofi/acfg/internal/configuration"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
	"testing"
)

func TestK8sManager_DeployOnEmptyCluster(t *testing.T) {
	const namespace = "default"
	const DeploymentName = "auth"
	kubeconfig := flag.String("kubeconfig", "/home/vahid/.kube/config", "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil{
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil{
		panic(err)
	}
	k := &K8s{clientset, namespace,[]string{DeploymentName}}
	fmt.Println(k)
	deploymentManager := clientset.AppsV1().Deployments(namespace)

	err = deploymentManager.Delete(context.TODO(), DeploymentName, v1.DeleteOptions{})
	if err != nil{
		if !k8serrors.IsNotFound(err){
			panic(err)
		}
	}
	fmt.Println("auth deployment removed!")

	fr, err := os.Open("/home/vahid/Desktop/projects/bookstore/k8s/auth-dep.yaml")
	if err != nil{
		panic(err)
	}

	err = k.Deploy(context.TODO(), fr)
	if err != nil{
		panic(err)
	}

	d, err := deploymentManager.Get(context.TODO(), DeploymentName, v1.GetOptions{})
	if err != nil{
		panic(err)
	}
	if d.Name != DeploymentName{
		fmt.Println(d.Name, "is not", DeploymentName)
		t.Fail()
		return
	}
	k.WaitAllDeploymentsAreStable(context.TODO())
	fmt.Println("TestK8sManager_DeployOnEmptyCluster", "DONE")
}

func TestK8sManager_UpdateWithNewConfig(t *testing.T) {
	TestK8sManager_DeployOnEmptyCluster(t)

	const namespace = "default"
	const DeploymentName = "auth"
	kubeconfig := "/home/vahid/.kube/config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil{
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil{
		panic(err)
	}
	k := &K8s{clientset, namespace,[]string{DeploymentName}}
	err = k.UpdateConfigurationsAndWait(context.TODO(), map[string]*configuration.Configuration{DeploymentName: &configuration.Configuration{
		ReplicaCount: int64Ptr(1),
		Memory: int64Ptr(128),
		CPU: int64Ptr(197),
	}})

	if err != nil{
		panic(err)
	}

	deploymentManager := clientset.AppsV1().Deployments(namespace)
	d, err := deploymentManager.Get(context.TODO(), DeploymentName, v1.GetOptions{})
	if err != nil{
		panic(err)
	}
	if d.Name != DeploymentName{
		t.Log(d.Name,"should be", DeploymentName)
		t.Fail()
	}
	if *d.Spec.Replicas != 1{
		t.Log("replica", *d.Spec.Replicas, "should be", 1)
		t.Fail()
	}
	fmt.Println(d.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
	cpuValue := d.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()

	if strings.Compare(cpuValue, "197m") != 0{
		t.Log("cpu", cpuValue, "should be", "197m", strings.Compare(cpuValue, "197m"))
		t.Fail()
	}

	k.WaitAllDeploymentsAreStable(context.TODO())
	fmt.Println("TestK8sManager_UpdateWithNewConfig", "DONE")
}

func TestK8sManager_UpdateWithTheSameConfig(t *testing.T) {
	TestK8sManager_UpdateWithNewConfig(t)

	const namespace = "default"
	const DeploymentName = "auth"
	kubeconfig := "/home/vahid/.kube/config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil{
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil{
		panic(err)
	}
	k := &K8s{clientset, namespace,[]string{DeploymentName}}

	err = k.UpdateConfigurationsAndWait(context.TODO(), map[string]*configuration.Configuration{DeploymentName: &configuration.Configuration{
		ReplicaCount: int64Ptr(1),
		Memory: int64Ptr(128),
		CPU: int64Ptr(197),
	}})

	if err != nil{
		panic(err)
	}
	k.WaitAllDeploymentsAreStable(context.TODO())
}