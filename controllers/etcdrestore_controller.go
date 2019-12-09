package controllers

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	etcdv1alpha1 "github.com/improbable-eng/etcd-cluster-operator/api/v1alpha1"
)

// EtcdRestoreReconciler reconciles a EtcdRestore object
type EtcdRestoreReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=etcd.improbable.io,resources=etcdrestores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=etcd.improbable.io,resources=etcdrestores/status,verbs=get;update;patch

func (r *EtcdRestoreReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("etcdrestore", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *EtcdRestoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&etcdv1alpha1.EtcdRestore{}).
		Complete(r)
}
