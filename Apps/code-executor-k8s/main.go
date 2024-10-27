package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func boolPtr(a bool) *bool {
	return &a
}

func int32Ptr(a int32) *int32 {
	return &a
}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	c := codeHandler{
		client:           clientset,
		workingNamespace: "default",
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		// pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		// if err != nil {
		// 	panic(err.Error())
		// }
		// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// podLogOpts := core.PodLogOptions{
		// 	Container: "alertmanager",
		// }
		// req := clientset.CoreV1().Pods("default").GetLogs("alertmanager-monitoring-kube-prometheus-alertmanager-0", &podLogOpts)
		// zz, err := req.Stream(context.TODO())
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// defer zz.Close()

		// buf := new(bytes.Buffer)
		// _, err = io.Copy(buf, zz)
		// str := buf.String()
		// fmt.Printf("Print logs for %v\n", "alertmanager-monitoring-kube-prometheus-alertmanager-0")
		// fmt.Println(str)
		// fmt.Printf("End collecting logs for %v\n", "alertmanager-monitoring-kube-prometheus-alertmanager-0")

		c.deleteJob("test-test")
		c.deleteConfigmap("test-test")
		c.createCodeConfigmap()
		jj := c.createJob()
		fmt.Printf("Completions: %v\n", *jj.Spec.Completions)
		time.Sleep(10 * time.Second)
		if *jj.Spec.Completions < 1 {
			fmt.Println("still waiting for job to complete")
			time.Sleep(5 * time.Second)
		}
		podName, err := c.getPodName("zzz=zzz")
		if err != nil {
			fmt.Printf("require further investigation: %v\n", err)
			continue
		}
		yahoo := c.getPodLogs(podName)
		fmt.Printf("Logs from pod: %v\n", yahoo)

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		// _, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
		// if errors.IsNotFound(err) {
		// 	fmt.Printf("Pod example-xxxxx not found in default namespace\n")
		// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		// } else if err != nil {
		// 	panic(err.Error())
		// } else {
		// 	fmt.Printf("Found example-xxxxx pod in default namespace\n")
		// }

		time.Sleep(10 * time.Second)
	}
}

type codeHandler struct {
	client           *kubernetes.Clientset
	workingNamespace string
}

func (c codeHandler) getPodName(labelFilter string) (string, error) {
	fmt.Printf("LabelFilter: %v\n", labelFilter)
	listOpts := metav1.ListOptions{
		LabelSelector: labelFilter,
	}
	r, err := c.client.CoreV1().Pods(c.workingNamespace).List(context.TODO(), listOpts)
	if err != nil {
		fmt.Println("unable to list out the required Pod")
	}
	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Second)
		if len(r.Items) < 1 {
			fmt.Printf("found %v pods\n", len(r.Items))
			continue
		}
		return r.Items[0].ObjectMeta.Name, nil
	}
	return "", fmt.Errorf("unable to find pods with labelfilter: %v", labelFilter)
}

func (c codeHandler) getPodLogs(podName string) string {
	podLogOpts := core.PodLogOptions{
		Container: "test",
	}
	req := c.client.CoreV1().Pods(c.workingNamespace).GetLogs(podName, &podLogOpts)
	zz, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	defer zz.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, zz)
	logs := buf.String()
	return logs
}

func (c codeHandler) createCodeConfigmap() {
	fmt.Println("start creating test-test config")
	cConfig := core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-test",
		},
		Immutable: boolPtr(true),
		Data: map[string]string{
			"code": "import time\ntime.sleep(30)\nprint(\"test\")",
		},
	}
	createOpts := metav1.CreateOptions{}
	_, err := c.client.CoreV1().ConfigMaps(c.workingNamespace).Create(context.TODO(), &cConfig, createOpts)
	if err != nil {
		fmt.Println(err)
	}
}

func (c codeHandler) deleteConfigmap(configmapName string) error {
	fmt.Printf("delete configmaps %v\n", configmapName)
	deleteOpts := metav1.DeleteOptions{}
	err := c.client.CoreV1().ConfigMaps(c.workingNamespace).Delete(context.TODO(), configmapName, deleteOpts)
	if err != nil {
		fmt.Println(err)
	}
	getOpts := metav1.GetOptions{}
	for i := 0; i < 10; i++ {
		_, err := c.client.CoreV1().ConfigMaps(c.workingNamespace).Get(context.TODO(), configmapName, getOpts)
		if err != nil {
			fmt.Printf("unable to get configmap: %v - highly likely\n", configmapName)
			return nil
		}
		fmt.Println("configmaps can still be found. we need to wait first")
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("configmaps can still be found - need to investigate")
}

// deleteJob in synchronous fashion
func (c codeHandler) deleteJob(jobName string) error {
	fmt.Printf("delete job %v\n", jobName)
	df := metav1.DeletePropagationForeground
	deleteOpts := metav1.DeleteOptions{
		PropagationPolicy: &df,
	}
	err := c.client.BatchV1().Jobs(c.workingNamespace).Delete(context.TODO(), jobName, deleteOpts)
	if err != nil {
		fmt.Println(err)
	}
	getOpts := metav1.GetOptions{}
	for i := 0; i < 10; i++ {
		_, err := c.client.BatchV1().Jobs(c.workingNamespace).Get(context.TODO(), jobName, getOpts)
		if err != nil {
			fmt.Printf("unable to get job: %v - highly likely\n", jobName)
			return nil
		}
		fmt.Println("job can still be found. we need to wait first")
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("job can still be found - need to investigate")
}

func (c codeHandler) createJob() *batchv1.Job {
	fmt.Println("start creating create job test-test")
	jConfig := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-test",
		},
		Spec: batchv1.JobSpec{

			Parallelism: int32Ptr(1),
			Completions: int32Ptr(1),
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"zzz": "zzz",
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name: "test",
							// Image: "python:3.9-alpine",
							Image:   "nginx:latest",
							Command: []string{"cat", "/code/code"},
							// Command: []string{"python", "/code/code"},
							VolumeMounts: []core.VolumeMount{
								{Name: "miao", ReadOnly: true, MountPath: "/code"},
							},
						},
					},
					Volumes: []core.Volume{
						{
							Name: "miao",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "test-test",
									},
								},
							},
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}
	createOpts := metav1.CreateOptions{}
	jj, err := c.client.BatchV1().Jobs(c.workingNamespace).Create(context.TODO(), &jConfig, createOpts)
	if err != nil {
		fmt.Println(err)
	}
	return jj
}
