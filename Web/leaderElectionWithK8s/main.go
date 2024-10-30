package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

var clientset *kubernetes.Clientset
var leaderState bool

func main() {
	time.Sleep(15 * time.Second)
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	POD_NAME := os.Getenv("POD_NAME")

	rl, err := resourcelock.NewFromKubeconfig(resourcelock.LeasesResourceLock, "default", "app-lock", resourcelock.ResourceLockConfig{
		Identity: POD_NAME,
	}, config, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		panic("Failed due to bad resource lock")
	}
	ctx := context.Background()
	LESettings := leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: 10 * time.Second,
		RenewDeadline: 5 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: zzz,
			OnStoppedLeading: func() {
				fmt.Println("Stopped")
				panic("stopped leading")
			},
			OnNewLeader: func(id string) {
				if id != POD_NAME {
					fmt.Println("is not the leader")
					leaderState = false
				} else {
					fmt.Println("is the leader")
					leaderState = true
				}
			},
		},
		Name: "debugging",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request on leader state")
		defer fmt.Println("End request on leader state")
		w.Write([]byte(fmt.Sprintf("LeaderState: %v", leaderState)))
	})
	go http.ListenAndServe(":8080", nil)

	leaderelection.RunOrDie(ctx, LESettings)

}

func zzz(ctx context.Context) {
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod example-xxxxx not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found example-xxxxx pod in default namespace\n")
		}

		time.Sleep(10 * time.Second)
	}
}
