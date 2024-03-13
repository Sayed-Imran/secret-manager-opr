/*
Copyright 2024 Sayed Imran.

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

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/Sayed-Imran/secret-manager-opr/api/v1alpha1"
)

// SecretManagerReconciler reconciles a SecretManager object
type SecretManagerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.secrets.mgr,resources=secretmanagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.secrets.mgr,resources=secretmanagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.secrets.mgr,resources=secretmanagers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SecretManager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *SecretManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	secretManager := &apiv1alpha1.SecretManager{}
	err := r.Get(ctx, req.NamespacedName, secretManager)
	if err != nil {
		log.Log.Error(err, "unable to fetch SecretManager")
		return ctrl.Result{}, err
	}
	allNamespaces := v1.NamespaceList{}
	err = r.List(ctx, &allNamespaces)
	if err != nil {
		log.Log.Error(err, "unable to fetch all namespaces")
		return ctrl.Result{}, err
	}


	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.SecretManager{}).
		Complete(r)
}

func getMatchedNamespaces(matchNamespaces []string, avoidNamespaces []string, ctx context.Context) []string {
	var matchedNamespaces []string
	allNamespaces := v1.NamespaceList{}
	err := r.List(ctx, &allNamespaces)
	if err != nil {
		log.Log.Error(err, "unable to fetch all namespaces")
		return matchedNamespaces
	}
	for _, namespace := range allNamespaces.Items {
		if len(matchNamespaces) > 0 {
			for _, namespace := range matchNamespaces {
				// implement regex check

			}
		} else {
			matchedNamespaces = append(matchedNamespaces, namespace.Name)
		}
	}
	return matchedNamespaces
}
