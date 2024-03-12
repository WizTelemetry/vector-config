/*
Copyright 2023.

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
	"github.com/kubesphere-sigs/vector-config/internal/constants"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// SecretReconciler reconciles a Secret object
type SecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups="",resources=secrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Secret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *SecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	configsecrets := &v1.SecretList{}
	configselector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			constants.SecretLabel:         constants.VectorRole,
			constants.ConfigReloadEnabled: "true",
		},
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	CAsecrets := &v1.SecretList{}
	CAselector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			constants.SecretLabel:        constants.VectorRole,
			constants.CertificationLabel: "true",
		},
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	// remove all
	err = r.removeAllFiles()
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.List(ctx, configsecrets, &client.ListOptions{
		Namespace:     req.Namespace,
		LabelSelector: configselector,
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	configpath := fmt.Sprintf("%s", constants.FileDir)
	err = os.MkdirAll(configpath, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return ctrl.Result{}, err
	}
	capath := fmt.Sprintf("%s/%s", constants.FileDir, constants.Certification)
	err = os.MkdirAll(capath, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return ctrl.Result{}, err
	}

	for _, cs := range configsecrets.Items {
		for s, bytes := range cs.Data {
			path := fmt.Sprintf("%s/%s-%s", constants.FileDir, cs.Name, s)
			err = os.WriteFile(path, bytes, 0644)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	err = r.List(ctx, CAsecrets, &client.ListOptions{
		Namespace:     req.Namespace,
		LabelSelector: CAselector,
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, cs := range CAsecrets.Items {
		for s, bytes := range cs.Data {
			path := fmt.Sprintf("%s/%s/%s", constants.FileDir, constants.Certification, s)
			err = os.WriteFile(path, bytes, 0644)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		For(&v1.Secret{}).
		Complete(r)
}

func (r *SecretReconciler) removeAllFiles() error {
	err := filepath.Walk(constants.FileDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
