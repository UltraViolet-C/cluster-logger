/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	logv1 "my.domain/clusterlogger/api/v1"
)

// ClusterScanReconciler reconciles a ClusterScan object
type ClusterScanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Clock
}

type realClock struct{}

func (_ realClock) Now() time.Time { return time.Now() }

// mocked out so timing can be faked during tests
type Clock interface {
	Now() time.Time
}

func prettyPrintClusterInfo(clusterScan logv1.ClusterScan) string {
	var output string = fmt.Sprintf(">>> Cluster Scan made at time %s <<<\n", time.Now().String()) +
		fmt.Sprintf("Version: %s | Name: %s\n", clusterScan.Spec.Version, clusterScan.Spec.Name) +
		"Nodes:\n"

	var nameColumnWidth int = len("Name")
	var uidColumnWidth int = len("UID")
	var podNumberColumnWidth int = len("# of Pods")

	// loop once to find out how much padding is required per column
	for i := 0; i < len(clusterScan.Spec.Nodes); i++ {
		nameColumnWidth = int(math.Max(float64(nameColumnWidth), float64(len(clusterScan.Spec.Nodes[i].Name))))
		// using Itoa is definitely not the fastest way to find this but I didn't want to add a loop
		// and a couple temp variables when I could do it with one line of code.
		uidColumnWidth = int(math.Max(float64(uidColumnWidth), float64(len(strconv.Itoa(int(clusterScan.Spec.Nodes[i].UID))))))
		podNumberColumnWidth = int(math.Max(float64(podNumberColumnWidth), float64(len(strconv.Itoa(int(clusterScan.Spec.Nodes[i].NumberOfPods)))+5)))
	}

	output = output + strings.Repeat(" ", 4) + "Name" + strings.Repeat(" ", nameColumnWidth-len("Name")) +
		" | UID" + strings.Repeat(" ", uidColumnWidth-len("UID")) +
		" | # of Pods" + strings.Repeat(" ", podNumberColumnWidth-len("# of Pods")) + " | Status\n"

	// another loop to actually add the data to output string
	for i := 0; i < len(clusterScan.Spec.Nodes); i++ {
		var line string = ""
		if clusterScan.Spec.Nodes[i].Master {
			line = line + "  * "
		} else {
			line = line + strings.Repeat(" ", 4)
		}
		line = line + clusterScan.Spec.Nodes[i].Name + strings.Repeat(" ", nameColumnWidth-len(clusterScan.Spec.Nodes[i].Name)) + " | " +
			fmt.Sprintf("%d", clusterScan.Spec.Nodes[i].UID) + strings.Repeat(" ", uidColumnWidth-len(strconv.Itoa(int(clusterScan.Spec.Nodes[i].UID)))) + " | " +
			fmt.Sprintf("%d Pods", clusterScan.Spec.Nodes[i].NumberOfPods) + strings.Repeat(" ", podNumberColumnWidth-len(fmt.Sprintf("%d Pods", clusterScan.Spec.Nodes[i].NumberOfPods))) + " | " +
			string(clusterScan.Spec.Nodes[i].Status)

		output = output + line + "\n"
	}

	return output
}

//+kubebuilder:rbac:groups=log.my.domain,resources=clusterscans,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=log.my.domain,resources=clusterscans/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=log.my.domain,resources=clusterscans/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterScan object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ClusterScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	// load the clusterScan object
	var clusterScan logv1.ClusterScan
	if err := r.Get(ctx, req.NamespacedName, &clusterScan); err != nil {
		log.Error(err, "unable to fetch ClusterScan")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// temp log just to know that this all works
	log.V(1).Info("loaded clusterScan object!")

	var cmdString string = prettyPrintClusterInfo(clusterScan)
	constructLogJob := func(clusterScan *logv1.ClusterScan) (*kbatch.Job, error) {
		name := fmt.Sprintf("clusterscanlog-%d", time.Now().Unix())

		// make a string here with the pretty printed representation of our nodes
		// and feed it into the Job's Command field and we're done! :D

		job := &kbatch.Job{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      make(map[string]string),
				Annotations: make(map[string]string),
				Name:        name,
				Namespace:   clusterScan.Namespace,
			},
			Spec: kbatch.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:    name,
								Image:   "bash",
								Command: strings.Split(fmt.Sprintf("echo %s", cmdString), " "),
							},
						},
						RestartPolicy: corev1.RestartPolicyNever,
					},
				},
			},
		}
		if err := ctrl.SetControllerReference(clusterScan, job, r.Scheme); err != nil {
			return nil, err
		}

		return job, nil
	}

	// construct our job
	job, err := constructLogJob(&clusterScan)
	if err != nil {
		log.Error(err, "unable to construct job from template")
		return ctrl.Result{}, nil
	}

	// create job on cluster
	if err := r.Create(ctx, job); err != nil {
		log.Error(err, "unable to create Job for CronJob", "job", job)
		return ctrl.Result{}, err
	}

	// log success at highest verbosity level for readability in logs
	log.V(1).Info("created Job for ClusterScan", "job", job)

	// return nothin if the job is correctly constructed and created on cluster
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// set up a real clock when not in a test
	if r.Clock == nil {
		r.Clock = realClock{}
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&logv1.ClusterScan{}).
		Complete(r)
}
