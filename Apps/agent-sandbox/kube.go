package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var (
	claimGVR    = schema.GroupVersionResource{Group: "extensions.agents.x-k8s.io", Version: "v1alpha1", Resource: "sandboxclaims"}
	sandboxGVR  = schema.GroupVersionResource{Group: "agents.x-k8s.io", Version: "v1alpha1", Resource: "sandboxes"}
	snapshotGVR = schema.GroupVersionResource{Group: "podsnapshot.gke.io", Version: "v1", Resource: "podsnapshots"}
	triggerGVR  = schema.GroupVersionResource{Group: "podsnapshot.gke.io", Version: "v1", Resource: "podsnapshotmanualtriggers"}
)

type kubeService struct {
	core       kubernetes.Interface
	dynamic    dynamic.Interface
	namespace  string
	templateID string
	warmPool   string
}

type dashboardData struct {
	Namespace string
	Template  string
	WarmPool  string
	Claims    []claimView
	Pods      []podView
	Snapshots []snapshotView
	Message   string
	Error     string
	UpdatedAt string
}

type claimView struct {
	Name        string
	Sandbox     string
	Ready       string
	Replicas    int64
	CreatedAt   string
	SnapshotNum int
}

type podView struct {
	Name       string
	Phase      string
	Ready      string
	Node       string
	IP         string
	Runtime    string
	Age        string
	Restored   string
	LastOutput string
}

type snapshotView struct {
	Name      string
	Pod       string
	Status    string
	CreatedAt string
	Storage   string
}

func newKubeService(core kubernetes.Interface, dynamicClient dynamic.Interface, namespace, templateID, warmPool string) *kubeService {
	return &kubeService{core: core, dynamic: dynamicClient, namespace: namespace, templateID: templateID, warmPool: warmPool}
}

func (k *kubeService) dashboard(ctx context.Context) (dashboardData, error) {
	claims, err := k.dynamic.Resource(claimGVR).Namespace(k.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return dashboardData{}, err
	}
	sandboxes, err := k.dynamic.Resource(sandboxGVR).Namespace(k.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return dashboardData{}, err
	}
	pods, err := k.core.CoreV1().Pods(k.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return dashboardData{}, err
	}
	snapshots, err := k.dynamic.Resource(snapshotGVR).Namespace(k.namespace).List(ctx, metav1.ListOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return dashboardData{}, err
	}

	sandboxByName := make(map[string]unstructured.Unstructured, len(sandboxes.Items))
	for _, item := range sandboxes.Items {
		sandboxByName[item.GetName()] = item
	}
	snapshotCounts := map[string]int{}
	snapshotViews := make([]snapshotView, 0, len(snapshots.Items))
	for _, item := range snapshots.Items {
		hash := item.GetLabels()["agents.x-k8s.io/sandbox-name-hash"]
		snapshotCounts[hash]++
		snapshotViews = append(snapshotViews, snapshotView{
			Name: item.GetName(), Pod: item.GetAnnotations()["podsnapshot.gke.io/origin-pod"],
			Status: snapshotStatus(&item), CreatedAt: age(item.GetCreationTimestamp().Time),
			Storage: nestedString(item.Object, "status", "storageStatus", "gcs", "observedGCSPath"),
		})
	}
	sort.Slice(snapshotViews, func(i, j int) bool { return snapshotViews[i].Name > snapshotViews[j].Name })

	claimViews := make([]claimView, 0, len(claims.Items))
	for _, claim := range claims.Items {
		sandboxName := nestedString(claim.Object, "status", "sandbox", "name")
		view := claimView{Name: claim.GetName(), Sandbox: sandboxName, Ready: condition(&claim, "Ready"), CreatedAt: age(claim.GetCreationTimestamp().Time)}
		if sandbox, ok := sandboxByName[sandboxName]; ok {
			view.Replicas = nestedInt64(sandbox.Object, "spec", "replicas")
			hash := selectorValue(nestedString(sandbox.Object, "status", "selector"), "agents.x-k8s.io/sandbox-name-hash")
			view.SnapshotNum = snapshotCounts[hash]
		}
		claimViews = append(claimViews, view)
	}
	sort.Slice(claimViews, func(i, j int) bool { return claimViews[i].Name < claimViews[j].Name })

	podViews := make([]podView, 0, len(pods.Items))
	for i := range pods.Items {
		pod := &pods.Items[i]
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		restored := "-"
		for _, c := range pod.Status.Conditions {
			if c.Type == "PodRestored" {
				restored = string(c.Status)
				if c.Message != "" {
					restored += ": " + c.Message
				}
			}
		}
		podViews = append(podViews, podView{
			Name: pod.Name, Phase: string(pod.Status.Phase), Ready: podReady(pod), Node: pod.Spec.NodeName,
			IP: pod.Status.PodIP, Runtime: valueOr(pod.Spec.RuntimeClassName, "default"),
			Age: age(pod.CreationTimestamp.Time), Restored: restored, LastOutput: k.lastLogLine(ctx, pod),
		})
	}
	sort.Slice(podViews, func(i, j int) bool { return podViews[i].Name < podViews[j].Name })

	return dashboardData{
		Namespace: k.namespace, Template: k.templateID, WarmPool: k.warmPool,
		Claims: claimViews, Pods: podViews, Snapshots: snapshotViews,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (k *kubeService) createClaim(ctx context.Context, requestedName string) (string, error) {
	if requestedName != "" {
		if errs := validation.IsDNS1123Subdomain(requestedName); len(errs) > 0 {
			return "", fmt.Errorf("invalid claim name: %s", strings.Join(errs, ", "))
		}
	}
	claim := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "extensions.agents.x-k8s.io/v1alpha1",
		"kind":       "SandboxClaim",
		"metadata": map[string]any{
			"namespace": k.namespace,
			"labels":    map[string]any{"app.kubernetes.io/managed-by": "agent-sandbox-demo"},
		},
		"spec": map[string]any{
			"sandboxTemplateRef": map[string]any{"name": k.templateID},
			"warmpool":           k.warmPool,
			"lifecycle":          map[string]any{"shutdownPolicy": "DeleteForeground"},
		},
	}}
	if requestedName == "" {
		claim.SetGenerateName("demo-")
	} else {
		claim.SetName(requestedName)
	}
	created, err := k.dynamic.Resource(claimGVR).Namespace(k.namespace).Create(ctx, claim, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return created.GetName(), nil
}

func (k *kubeService) destroyClaim(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("claim name is required")
	}
	propagation := metav1.DeletePropagationForeground
	return k.dynamic.Resource(claimGVR).Namespace(k.namespace).Delete(ctx, name, metav1.DeleteOptions{PropagationPolicy: &propagation})
}

func (k *kubeService) setReplicas(ctx context.Context, sandbox string, replicas int64) error {
	if sandbox == "" {
		return fmt.Errorf("sandbox name is required")
	}
	body, _ := json.Marshal(map[string]any{"spec": map[string]any{"replicas": replicas}})
	_, err := k.dynamic.Resource(sandboxGVR).Namespace(k.namespace).Patch(ctx, sandbox, types.MergePatchType, body, metav1.PatchOptions{})
	return err
}

func (k *kubeService) snapshot(ctx context.Context, sandboxName string) (string, error) {
	podName, _, err := k.sandboxPod(ctx, sandboxName)
	if err != nil {
		return "", err
	}
	trigger := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "podsnapshot.gke.io/v1", "kind": "PodSnapshotManualTrigger",
		"metadata": map[string]any{"generateName": safePrefix(sandboxName) + "-ui-", "namespace": k.namespace},
		"spec":     map[string]any{"targetPod": podName},
	}}
	created, err := k.dynamic.Resource(triggerGVR).Namespace(k.namespace).Create(ctx, trigger, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		current, err := k.dynamic.Resource(triggerGVR).Namespace(k.namespace).Get(ctx, created.GetName(), metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		for _, c := range nestedConditions(current.Object) {
			if c["type"] == "Triggered" && c["status"] == "True" && c["reason"] == "Complete" {
				name := nestedString(current.Object, "status", "snapshotCreated", "name")
				if name == "" {
					return "", fmt.Errorf("trigger completed without snapshot name")
				}
				return name, nil
			}
			if c["type"] == "Triggered" && c["status"] == "False" {
				return "", fmt.Errorf("snapshot failed: %s", c["message"])
			}
		}
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("wait for snapshot: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func (k *kubeService) restoreLatest(ctx context.Context, sandboxName string) (string, error) {
	podName, hash, podErr := k.sandboxPod(ctx, sandboxName)
	if podErr != nil && !apierrors.IsNotFound(podErr) {
		return "", podErr
	}
	if apierrors.IsNotFound(podErr) {
		podName = ""
	}
	if hash == "" {
		sandbox, err := k.dynamic.Resource(sandboxGVR).Namespace(k.namespace).Get(ctx, sandboxName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		hash = selectorValue(nestedString(sandbox.Object, "status", "selector"), "agents.x-k8s.io/sandbox-name-hash")
	}
	if hash == "" {
		return "", fmt.Errorf("sandbox has no snapshot grouping label")
	}
	list, err := k.dynamic.Resource(snapshotGVR).Namespace(k.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set{"agents.x-k8s.io/sandbox-name-hash": hash}.String(),
	})
	if err != nil {
		return "", err
	}
	var latest *unstructured.Unstructured
	for i := range list.Items {
		item := &list.Items[i]
		if !conditionTrue(item, "Ready") {
			continue
		}
		if latest == nil || item.GetCreationTimestamp().After(latest.GetCreationTimestamp().Time) {
			latest = item
		}
	}
	if latest == nil {
		return "", fmt.Errorf("no ready snapshot exists for sandbox %s", sandboxName)
	}
	if podName != "" {
		grace := int64(0)
		if err := k.core.CoreV1().Pods(k.namespace).Delete(ctx, podName, metav1.DeleteOptions{GracePeriodSeconds: &grace}); err != nil && !apierrors.IsNotFound(err) {
			return "", err
		}
	} else if err := k.setReplicas(ctx, sandboxName, 1); err != nil {
		return "", err
	}
	return latest.GetName(), nil
}

func (k *kubeService) sandboxPod(ctx context.Context, sandboxName string) (string, string, error) {
	sandbox, err := k.dynamic.Resource(sandboxGVR).Namespace(k.namespace).Get(ctx, sandboxName, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}
	hash := selectorValue(nestedString(sandbox.Object, "status", "selector"), "agents.x-k8s.io/sandbox-name-hash")
	podName := sandbox.GetAnnotations()["agents.x-k8s.io/pod-name"]
	if podName == "" {
		podName = sandboxName
	}
	_, err = k.core.CoreV1().Pods(k.namespace).Get(ctx, podName, metav1.GetOptions{})
	return podName, hash, err
}

func (k *kubeService) lastLogLine(ctx context.Context, pod *corev1.Pod) string {
	if len(pod.Spec.Containers) == 0 || pod.Status.Phase != corev1.PodRunning {
		return "-"
	}
	tail := int64(1)
	stream, err := k.core.CoreV1().Pods(k.namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
		Container: pod.Spec.Containers[0].Name, TailLines: &tail,
	}).Stream(ctx)
	if err != nil {
		return "-"
	}
	defer stream.Close()
	data, err := io.ReadAll(io.LimitReader(stream, 4096))
	if err != nil {
		return "-"
	}
	return strings.TrimSpace(string(data))
}

func condition(obj *unstructured.Unstructured, name string) string {
	for _, c := range nestedConditions(obj.Object) {
		if c["type"] == name {
			return c["status"]
		}
	}
	return "Unknown"
}

func conditionTrue(obj *unstructured.Unstructured, name string) bool {
	return condition(obj, name) == "True"
}

func snapshotStatus(obj *unstructured.Unstructured) string {
	for _, c := range nestedConditions(obj.Object) {
		if c["type"] == "Ready" {
			if c["status"] == "True" {
				return c["reason"]
			}
			return c["status"] + ": " + c["reason"]
		}
	}
	return "Pending"
}

func nestedConditions(object map[string]any) []map[string]string {
	raw, _, _ := unstructured.NestedSlice(object, "status", "conditions")
	out := make([]map[string]string, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		out = append(out, map[string]string{
			"type": fmt.Sprint(m["type"]), "status": fmt.Sprint(m["status"]),
			"reason": fmt.Sprint(m["reason"]), "message": fmt.Sprint(m["message"]),
		})
	}
	return out
}

func nestedString(object map[string]any, fields ...string) string {
	value, _, _ := unstructured.NestedString(object, fields...)
	return value
}

func nestedInt64(object map[string]any, fields ...string) int64 {
	value, _, _ := unstructured.NestedInt64(object, fields...)
	return value
}

func selectorValue(selector, key string) string {
	for _, term := range strings.Split(selector, ",") {
		parts := strings.SplitN(term, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			return parts[1]
		}
	}
	return ""
}

func podReady(pod *corev1.Pod) string {
	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.PodReady {
			return string(c.Status)
		}
	}
	return "Unknown"
}

func valueOr(value *string, fallback string) string {
	if value != nil {
		return *value
	}
	return fallback
}

func age(created time.Time) string {
	if created.IsZero() {
		return "-"
	}
	d := time.Since(created).Round(time.Second)
	if d < time.Minute {
		return d.String()
	}
	if d < time.Hour {
		return d.Round(time.Minute).String()
	}
	return d.Round(time.Hour).String()
}

func safePrefix(value string) string {
	value = strings.Trim(value, "-")
	if len(value) > 40 {
		value = value[:40]
	}
	if value == "" {
		return "sandbox"
	}
	return value
}
