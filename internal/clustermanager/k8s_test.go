package clustermanager

//import (
//	"context"
//	"flag"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/tools/clientcmd"
//	"os"
//	"testing"
//)
//
//func TestK8sManager_Deploy(t *testing.T) {
//	kubeconfig := flag.String("kubeconfig", "/home/vahid/.kube/config", "kubeconfig file")
//	flag.Parse()
//	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
//	if err != nil{
//		panic(err)
//	}
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil{
//		panic(err)
//	}
//
//	k := &K8sManager{clientset, "default"}
//
//	fr, err := os.Open("/home/vahid/Desktop/projects/bookstore/k8s/auth-dep.yaml")
//	if err != nil{
//		panic(err)
//	}
//
//	err = k.Deploy(context.TODO(), fr)
//	if err != nil{
//		panic(err)
//	}
//}
