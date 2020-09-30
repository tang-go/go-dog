package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", filepath.Join("./config", "kube.config"), "(optional) absolute path to the kubeconfig file")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get nodes:", err)
	}

	for i, node := range nodes.Items {
		fmt.Printf("[%d] %s\n", i, node.GetName())
		fmt.Println(node.Name)
		fmt.Println(node.CreationTimestamp) //加入集群时间
		fmt.Println(node.Status.NodeInfo)
		fmt.Println(node.Status.Conditions[len(node.Status.Conditions)-1].Type)
		fmt.Println(node.Status.Allocatable.Memory().String())
	}

	pods, err := clientset.CoreV1().Pods("qa").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %v pods in the cluster\n", len(pods.Items))
	for _, pod := range pods.Items {

		fmt.Println("名称:", pod.Name)
		fmt.Println("创建时间", pod.CreationTimestamp) //创建时间
		fmt.Println("标签:", pod.Labels)
		fmt.Println("命名空间:", pod.Namespace)
		fmt.Println("node地址:", pod.Status.HostIP)
		fmt.Println("pod地址:", pod.Status.PodIP)
		fmt.Println("pod启动时间:", pod.Status.StartTime)
		fmt.Println("pod状态:", pod.Status)

		fmt.Println("pod状态:", pod.Status.Phase)
		fmt.Println("重启次数:", pod.Status.ContainerStatuses[0].RestartCount) //重启次数
		fmt.Println("重启时间:", pod.Status.ContainerStatuses[0].Image)        //获取重启时间
		break
	}
	// for {
	// 	pods, err := clientset.CoreV1().Pods("qa").List(context.TODO(), metav1.ListOptions{})
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	fmt.Printf("There are %v pods in the cluster\n", pods.Items[0])

	// 	// Examples for error handling:
	// 	// - Use helper functions like e.g. errors.IsNotFound()
	// 	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	// 	namespace := "qa"
	// 	pod := "gateway-hall"
	// 	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
	// 	if errors.IsNotFound(err) {
	// 		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	// 	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	// 		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
	// 			pod, namespace, statusError.ErrStatus.Message)
	// 	} else if err != nil {
	// 		panic(err.Error())
	// 	} else {
	// 		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	// 	}
	// 	time.Sleep(10 * time.Second)
	// }
}
