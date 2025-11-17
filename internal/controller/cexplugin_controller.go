/*
Copyright 2025.

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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cexv1alpha1 "github.com/Ankana252/CexPluginOperator/api/v1alpha1"
)

// CexPluginReconciler reconciles a CexPlugin object
type CexPluginReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cex.cex.dev,resources=cexplugins,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cex.cex.dev,resources=cexplugins/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cex.cex.dev,resources=cexplugins/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CexPlugin object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *CexPluginReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// TODO(user): your logic here
	var cm corev1.ConfigMap
	if err := r.Get(ctx, types.NamespacedName{Name: "cex-config",
		Namespace: "cex-device-plugin"}, &cm); err != nil {
		// Handle error
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.

func (r *CexPluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cexv1alpha1.CexPlugin{}). // Primary resource
		Watches(
			&corev1.ConfigMap{}, // Secondary resource
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				// Always reconcile the same CR (or map logic if multiple CRs exist)
				return []reconcile.Request{
					{
						NamespacedName: types.NamespacedName{
							Name:      "cexplugin-sample",    // Your CR name
							Namespace: "cex-operator-system", // Use namespace of ConfigMap
						},
					},
				}
			}),
		).
		Complete(r)
}
