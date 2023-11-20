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
	"github.com/kubesphere-sigs/config-reload/internal/constants"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
	secret := &v1.Secret{}
	err := r.Get(ctx, req.NamespacedName, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// 其他错误
		return ctrl.Result{}, err
	}

	if !secret.ObjectMeta.DeletionTimestamp.IsZero() {
		// 检查是否包含Finalizer
		containsFinalizer := controllerutil.ContainsFinalizer(secret, constants.SecretFinalizer)
		if containsFinalizer {
			// 删除文件
			err = r.deleteFiles(ctx, secret)
			if err != nil {
				return ctrl.Result{}, err
			}
			// 移除Finalizer
			controllerutil.RemoveFinalizer(secret, constants.SecretFinalizer)
			err = r.Update(ctx, secret)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	} else {
		// 检查是否包含Finalizer
		containsFinalizer := controllerutil.ContainsFinalizer(secret, constants.SecretFinalizer)
		if !containsFinalizer {
			// 添加Finalizer
			controllerutil.AddFinalizer(secret, constants.SecretFinalizer)
			err = r.Update(ctx, secret)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// 创建或更新文件
	for s, bytes := range secret.Data {
		// 写入文件
		path := fmt.Sprintf("%s/%s", constants.FilePath, s)
		err = os.WriteFile(path, bytes, 0644)
		if err != nil {
			return ctrl.Result{}, err
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

func (r *SecretReconciler) deleteFiles(ctx context.Context, secret *v1.Secret) error {
	logger := log.FromContext(ctx)
	// 执行文件清理操作
	for s := range secret.Data {
		// 删除文件
		path := fmt.Sprintf("%s/%s", constants.FilePath, s)
		err := os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			// 如果出现错误并且不是文件不存在的错误，返回错误
			return err
		}
		logger.Info("delete file", "path", path)
	}
	return nil
}
