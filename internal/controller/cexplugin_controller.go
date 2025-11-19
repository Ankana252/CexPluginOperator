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
	logger := logf.FromContext(ctx)
	logger.Info("Reconciling CexPlugin ConfigMap", "Name", req.Name, "Namespace", req.Namespace)
	// TODO(user): your logic here

	// Fetch the ConfigMap
	configMap := &corev1.ConfigMap{}
	err := r.Get(ctx, req.NamespacedName, configMap)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			logger.Info("ConfigMap not found, may have been deleted", "name", req.Name)
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	// Log the ConfigMap data
	logger.Info("ConfigMap change detected!",
		"name", configMap.Name,
		"namespace", configMap.Namespace,
		"data", configMap.Data)

	if data, ok := configMap.Data["cex_resources.json"]; ok {
		logger.Info("cex_resources.json content", "content", data)
	} else {
		logger.Info("cex_resources.json key not found in ConfigMap")
	}

	return ctrl.Result{}, nil
}

// findConfigMapsForReconcile maps ConfigMap events to reconcile requests
func (r *CexPluginReconciler) findConfigMapsForReconcile(ctx context.Context, configMap client.Object) []reconcile.Request {
	logger := logf.FromContext(ctx)

	// Only watch ConfigMaps named "cex-resources" in namespace "cex-device-plugin"
	if configMap.GetName() == "cex-resources" && configMap.GetNamespace() == "cex-device-plugin" {
		logger.Info("ConfigMap watch triggered", "name", configMap.GetName())

		// Return a reconcile request
		// You can return multiple requests if needed
		return []reconcile.Request{
			{
				NamespacedName: types.NamespacedName{
					Name:      configMap.GetName(),
					Namespace: configMap.GetNamespace(),
				},
			},
		}
	}

	return []reconcile.Request{}
}

func (r *CexPluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cexv1alpha1.CexPlugin{}).
		Watches(
			&corev1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.findConfigMapsForReconcile),
		).
		Complete(r)
}
