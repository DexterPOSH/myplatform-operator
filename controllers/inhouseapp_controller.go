/*
Copyright 2021.

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

package controllers

import (
	"context"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myplatformv1alpha1 "github.com/DexterPOSH/myplatform-operator/api/v1alpha1"
)

// InhouseAppReconciler reconciles a InhouseApp object
type InhouseAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myplatform.dexterposh.github.io,resources=inhouseapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myplatform.dexterposh.github.io,resources=inhouseapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myplatform.dexterposh.github.io,resources=inhouseapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the InhouseApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *InhouseAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here
	logger := log.Log.WithValues("inhouseApp", req.NamespacedName)

	logger.Info("InhouseApp Reconcile method...")

	// fetch the inhouseApp CR instance
	inhouseApp := &myplatformv1alpha1.InhouseApp{}
	err := r.Get(ctx, req.NamespacedName, inhouseApp)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("InhouseApp resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get InhouseApp instance")
		return ctrl.Result{}, err
	}

	// check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: inhouseApp.Name, Namespace: inhouseApp.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// define a new deployment
		dep := r.deploymentForInhouseApp(inhouseApp) // deploymentForInhouseApp() returns a deployme
		logger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			logger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// deployment created, return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get Deployment")
		// Reconcile failed due to error - requeue
		return ctrl.Result{}, err
	}

	// This point, we have the deployment object created
	// Ensure the deployment size is same as the spec
	replicas := inhouseApp.Spec.Replicas
	if *found.Spec.Replicas != replicas {
		found.Spec.Replicas = &replicas
		err = r.Update(ctx, found)
		if err != nil {
			logger.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated return and requeue
		// Requeue for any reason other than an error
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the InhouseApp status with pod names
	// List the pods for this InhouseApp's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(inhouseApp.Namespace),
		client.MatchingLabels(inhouseApp.GetLabels()),
	}

	if err = r.List(ctx, podList, listOpts...); err != nil {
		logger.Error(err, "Falied to list pods", "InhouseApp.Namespace", inhouseApp.Namespace, "InhouseApp.Name", inhouseApp.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Pods if needed
	if !reflect.DeepEqual(podNames, inhouseApp.Status.Pods) {
		inhouseApp.Status.Pods = podNames
		err := r.Update(ctx, inhouseApp)
		if err != nil {
			logger.Error(err, "Failed to update InhouseApp status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *InhouseAppReconciler) deploymentForInhouseApp(m *myplatformv1alpha1.InhouseApp) *appsv1.Deployment {
	ls := m.GetLabels()
	replicas := m.Spec.Replicas

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "dexterposh/myappp-dev", //hard-coded here, make this dynamic
						Name:  "inhouseAppDeployment",  //hard-coded here, make this dynamic
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "http",
						}},
					}},
				},
			},
		},
	}
	ctrl.SetControllerReference(m, deploy, r.Scheme)
	return deploy
}

// Utility function to iterate over pods and return the names slice
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// SetupWithManager sets up the controller with the Manager.
func (r *InhouseAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myplatformv1alpha1.InhouseApp{}).
		Complete(r)
}
