package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
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
	_ = codeHandler{
		client:           clientset,
		workingNamespace: "default",
	}
	// for {
	// 	c.deleteJob("test-test")
	// 	c.deleteConfigmap("test-test")
	// 	c.createCodeConfigmap("test-test", "import time\ntime.sleep(30)\nprint(\"test\")")
	// 	c.createJob("test-test", "test-test")
	// 	podName, err := c.getPodName("zzz=zzz")
	// 	if err != nil {
	// 		fmt.Printf("require further investigation: %v\n", err)
	// 		continue
	// 	}
	// 	yahoo := c.getPodLogs(podName)
	// 	fmt.Printf("Logs from pod: %v\n", yahoo)
	// }
	items := make(map[string]codeRecord)
	zzz := cdb{Items: items}
	r := mux.NewRouter()

	r.Handle("/status", status{}).Methods(http.MethodGet)
	r.Handle("/submit-code-page", submitCodePage{}).Methods(http.MethodGet)
	r.Handle("/submit-code", submitCode{zz: zzz}).Methods(http.MethodPost)
	r.Handle("/list-code", listCode{zz: zzz}).Methods(http.MethodGet, http.MethodPost)
	r.Handle("/get-code/{uid}", getCode{zz: zzz}).Methods(http.MethodGet)
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

type status struct{}

func (s status) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("status: status checker")
}

type codeDB interface {
	Store(id, language, code, status string, submittedTime time.Time)
	Update(id, status, logs string, completedTime time.Time)
	Get(id string) (language, code, status, logs string, submittedTime, completedTime time.Time)
}

type codeRecord struct {
	Language      string
	Code          string
	Status        string
	Logs          string
	SubmittedTime time.Time
	CompletedTime time.Time
}

type cdb struct {
	Items map[string]codeRecord
}

func (c cdb) Store(id, language, code, status string, submittedTime time.Time) {
	c.Items[id] = codeRecord{
		Language:      language,
		Code:          code,
		Status:        status,
		SubmittedTime: submittedTime,
	}
}
func (c cdb) Update(id, status, logs string, completedTime time.Time) {
	a := c.Items[id]
	a.Status = status
	a.Logs = logs
	a.CompletedTime = completedTime
	c.Items[id] = a
}
func (c cdb) Get(id string) (language, code, status, logs string, submittedTime, completedTime time.Time) {
	a := c.Items[id]
	return a.Language, a.Code, a.Status, a.Logs, a.SubmittedTime, a.CompletedTime
}

//go:embed templates/*.html
var viewsFS embed.FS

type submitCodePage struct{}

func (s submitCodePage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("submit code page started")
	t, err := template.ParseFS(viewsFS, "templates/aaa.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
}

type submitCode struct {
	zz cdb
}

func (s submitCode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("submit code")
	uid, _ := uuid.NewV4()
	language := r.FormValue("language")
	code := r.FormValue("code")
	status := "submitted"
	submittedTime := time.Now()
	s.zz.Store(uid.String(), language, code, status, submittedTime)
	fmt.Printf("ID record created: %v\n", uid.String())
	http.Redirect(w, r, "/list-code", http.StatusTemporaryRedirect)
}

type listCode struct {
	zz cdb
}

func (l listCode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	aa := ""
	for _, z := range l.zz.Items {
		aa = aa + fmt.Sprintf("%v\n", z)
	}
	w.Write([]byte(aa))
}

type getCode struct {
	zz cdb
}

func (g getCode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("id")
	singleRecord := g.zz.Items[uid]
	aa := map[string]string{
		"code":          singleRecord.Code,
		"language":      singleRecord.Language,
		"status":        singleRecord.Status,
		"logs":          singleRecord.Logs,
		"submittedTime": singleRecord.SubmittedTime.Format(time.RFC3339),
		"compltedTime":  singleRecord.CompletedTime.Format(time.RFC3339),
	}
	newRaw, err := json.Marshal(aa)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("bad status: %v", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(newRaw)
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

func (c codeHandler) createCodeConfigmap(configmapName, code string) {
	fmt.Println("start creating test-test config")
	cConfig := core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: configmapName,
		},
		Immutable: boolPtr(true),
		Data: map[string]string{
			"code": code,
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

func (c codeHandler) createJob(jobName, configmapName string) *batchv1.Job {
	fmt.Println("start creating create job test-test")
	jConfig := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
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
							Name:  "test",
							Image: "python:3.9-alpine",
							// Image:   "nginx:latest",
							// Command: []string{"cat", "/code/code"},
							Command: []string{"python", "/code/code"},
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
										Name: configmapName,
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

	getOpts := metav1.GetOptions{}
	for i := 0; i < 20; i++ {
		jj, err := c.client.BatchV1().Jobs(c.workingNamespace).Get(context.TODO(), jobName, getOpts)
		if err != nil {
			fmt.Println("error in getting value for job")
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Printf("succeeded: %v :: failed: %v\n", jj.Status.Succeeded, jj.Status.Failed)
		if (jj.Status.Succeeded + jj.Status.Failed) < 1 {
			fmt.Println("still waiting for job to complete")
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println("a condition was hit. we wil exit here")
			break
		}
	}

	return jj
}
