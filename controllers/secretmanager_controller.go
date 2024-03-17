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
	"reflect"
	"regexp"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/Sayed-Imran/secret-manager-opr/api/v1alpha1"
	"github.com/Sayed-Imran/secret-manager-opr/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		if client.IgnoreNotFound(err) == nil {
			return ctrl.Result{}, nil
		}
	}
	matchNamespaces := getMatchedNamespaces(r, secretManager.Spec.MatchNamespaces, secretManager.Spec.AvoidNamespaces, ctx)
	for _, namespace := range matchNamespaces {
		err := createSecrets(r, secretManager.Spec.Data, string(secretManager.Spec.Type), namespace, secretManager.Name, ctx)
		if err != nil {
			log.Log.Error(err, "unable to create secret")
			return ctrl.Result{}, err
		}

	}
	// update the status of the secret manager
	secretManager.Status.Namespaces = matchNamespaces
	err = r.Status().Update(ctx, secretManager)
	if err != nil {
		log.Log.Error(err, "unable to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Duration(10 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.SecretManager{}).
		Complete(r)
}

func getMatchedNamespaces(r *SecretManagerReconciler, matchNamespaces []string, avoidNamespaces []string, ctx context.Context) []string {
	var matchedNamespaces []string
	allNamespaces := v1.NamespaceList{}
	err := r.List(ctx, &allNamespaces)
	if err != nil {
		log.Log.Error(err, "unable to fetch all namespaces")
		return matchedNamespaces
	}
	for _, namespace := range matchNamespaces {
		for _, ns := range allNamespaces.Items {
			matched, err := regexp.MatchString(namespace, ns.Name)
			if err != nil {
				log.Log.Error(err, "unable to match namespace")
				continue
			}
			if matched {
				matchedNamespaces = append(matchedNamespaces, ns.Name)
			}
		}
	}
	var finalNamespaces []string
	for _, ns := range matchedNamespaces {
		avoid := false
		for _, namespace := range avoidNamespaces {
			matched, err := regexp.MatchString(namespace, ns)
			if err != nil {
				log.Log.Error(err, "unable to match namespace")
				continue
			}
			if matched {
				avoid = true
				break
			}
		}
		if !avoid {
			finalNamespaces = append(finalNamespaces, ns)
		}
	}
	return finalNamespaces
}

func createSecrets(conciler *SecretManagerReconciler, data map[string][]byte, secretType string, namespace string, secretName string, ctx context.Context) error {

	newSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: namespace},
		Data:       data,
		Type:       constants.SecretTypes[secretType],
	}

	existingSecret := &v1.Secret{}
	err := conciler.Get(ctx, client.ObjectKey{Namespace: namespace, Name: secretName}, existingSecret)
	if err != nil {
		err = conciler.Create(ctx, newSecret)
		if err != nil {
			log.Log.Error(err, "unable to create secret")
			return err
		}
		return nil
	}

	if !reflect.DeepEqual(existingSecret.Data, newSecret.Data) {

		existingSecret.Data = newSecret.Data
		err = conciler.Update(ctx, existingSecret)
		if err != nil {
			log.Log.Error(err, "unable to update secret")
			return err
		}
	} else {
		log.Log.Info("secret already exists and data is the same, no update needed")
	}

	return nil
}
