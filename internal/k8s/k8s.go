package k8s

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

func Run(){
	const namespace = "default"
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

	//pod, err := clientset.CoreV1().Pods("default").Get(context.Background(),"gateway-6f46df854-plqv9", metav1.GetOptions{})
	fmt.Println("DeploymentSets:")

	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	deployments, err := deploymentsClient.List(context.Background(), metav1.ListOptions{})
	if err != nil{
		panic(err)
	}
	targetedDeployments := make(map[string]*v1.Deployment, 0)
	for _,d := range deployments.Items{
		if name, ok := d.Spec.Template.Labels["app"]; ok {
			targetedDeployments[name] = d.DeepCopy()
		}
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(context.TODO(), targetedDeployments["auth"].Name, metav1.GetOptions{})
		if getErr != nil{
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}
		resultB, _ := result.Marshal()
		beforHashed := md5.Sum(resultB)

		result.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				"cpu":    resource.MustParse("563m"),
				"memory": resource.MustParse("512Mi"),
			},
			Requests: corev1.ResourceList{
				"cpu":    resource.MustParse("485m"),
				"memory": resource.MustParse("512Mi"),
			},
		}

		returned, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		returnedB, _ := returned.Marshal()
		currentHashed := md5.Sum(returnedB)

		isSame := true
		for i := range beforHashed{
			isSame = isSame && beforHashed[i] == currentHashed[i]
		}
		fmt.Println("Compared them:", isSame)
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")

	something, err := deploymentsClient.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	//time.Sleep(10 * time.Second)
	for e := range something.ResultChan(){
		fmt.Println(e.Type)
		if e.Type != "MODIFIED"{
			continue
		}

		fmt.Printf("Listing deployments in namespace %q:\n", namespace)
		list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		allDone := true
		for _, d := range list.Items {
			//fmt.Printf(" * %s (%d replicas) %d\n", d.Name, *d.Spec.Replicas, d.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().MilliValue())
			fmt.Println(d.Name, d.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().MilliValue(), d.Status.Replicas, d.Status.AvailableReplicas, d.Status.UnavailableReplicas, d.Status.ReadyReplicas, d.Status.UpdatedReplicas)
			isReady := (d.Status.Replicas == d.Status.ReadyReplicas) && (d.Status.UnavailableReplicas == 0)
			if !isReady{
				fmt.Println(d.Name, "is not ready")
			}
			allDone = allDone && isReady
		}
		fmt.Println("------")
		if allDone{
			something.Stop()
		}
	}

	//something, err = clientset.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//for e := range something.ResultChan(){
	//	fmt.Println(e.Type)
	//	fmt.Printf("Listing deployments in namespace %q:\n", namespace)
	//	list, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	//	if err != nil {
	//		panic(err)
	//	}
	//	allDone := true
	//	for _, d := range list.Items {
	//		allRunning := true
	//		for _,cs := range d.Status.ContainerStatuses{
	//			allRunning = allRunning && cs.State.Running != nil
	//			//fmt.Println(d.Name, cs.State.Running, d.Status.String())
	//			//fmt.Println(d.Name, cs.State.Terminated, d.Status.String())
	//			//fmt.Println(d.Name, cs.State.Waiting, d.Status.String())
	//		}
	//		if !allRunning{
	//			fmt.Println(d.Name, "has something not running...")
	//		}
	//		allDone = allDone && allRunning
	//		//isReady := (d.Status.Replicas == d.Status.ReadyReplicas) && (d.Status.UnavailableReplicas == 0)
	//		//if !isReady{
	//		//	fmt.Println(d.Name, "is not ready")
	//		//}
	//		//allDone = allDone && isReady
	//		//allDone := true
	//	}
	//	fmt.Println("------")
	//	if allDone{
	//		something.Stop()
	//	}
	//}
}


