package controllers

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	etcdv1alpha1 "github.com/improbable-eng/etcd-cluster-operator/api/v1alpha1"
)

// EtcdRestoreReconciler reconciles a EtcdRestore object
type EtcdRestoreReconciler struct {
	client.Client
	Log             logr.Logger
	RestorePodImage string
}

// +kubebuilder:rbac:groups=etcd.improbable.io,resources=etcdrestores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=etcd.improbable.io,resources=etcdrestores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=persistantvolumeclaims/status,verbs=get;update;patch;create
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;create

func name(o metav1.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: o.GetNamespace(),
		Name:      o.GetName(),
	}
}

const (
	restoreContainerName = "etcd-restore"
	etcdRestoreImage     = "FIXME"
	restoredFromLabel    = "etcd.improbable.io/restored-from"
	restorePodLabel      = "etcd.improbable.io/restore-pod"
)

func markPVC(restore etcdv1alpha1.EtcdRestore, pvc *corev1.PersistentVolumeClaim) {
	pvc.Labels[restoredFromLabel] = restore.Name
}

func IsOurPVC(restore etcdv1alpha1.EtcdRestore, pvc corev1.PersistentVolumeClaim) bool {
	return pvc.Labels[restoredFromLabel] == restore.Name
}

func markPod(restore etcdv1alpha1.EtcdRestore, pod *corev1.Pod) {
	pod.Labels[restoredFromLabel] = restore.Name
	pod.Labels[restorePodLabel] = "true"
}

func IsOurPod(restore etcdv1alpha1.EtcdRestore, pod corev1.Pod) bool {
	return pod.Labels[restoredFromLabel] == restore.Name &&
		pod.Labels[restorePodLabel] == "true"
}

func markCluster(restore etcdv1alpha1.EtcdRestore, cluster *etcdv1alpha1.EtcdCluster) {
	cluster.Labels[restoredFromLabel] = restore.Name
}

func IsOurCluster(restore etcdv1alpha1.EtcdRestore, cluster *etcdv1alpha1.EtcdCluster) bool {
	return cluster.Labels[restoredFromLabel] == restore.Name
}

func (r *EtcdRestoreReconciler) podForRestore(restore etcdv1alpha1.EtcdRestore, pvc *corev1.PersistentVolumeClaim) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      restore.Name,
			Namespace: restore.Namespace,
			Labels: map[string]string{
				restoredFromLabel: restore.Name,
				restorePodLabel:   "true",
			},
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "etcd-data-directory",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvc.Name,
							// We definitely need to be able to write to it.
							ReadOnly: false,
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:         "etcd-restore",
					Image:        r.RestorePodImage,
					Command:      nil,
					Args:         nil,
					VolumeMounts: nil,
				},
			},
			// If the Pod fails we should *not* attempt to recover automatically. Leave it failed.
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
}

func (r *EtcdRestoreReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("etcdrestore", req.NamespacedName)

	var restore etcdv1alpha1.EtcdRestore
	if err := r.Get(ctx, req.NamespacedName, &restore); err != nil {
		// Can't find the resource. We could have just been deleted? If so, no need to do anything as the owner
		// references will clean everything up with a cascading delete.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else {
		// Found the resource, continue with reconciliation.
	}

	if restore.Status.Phase == etcdv1alpha1.EtcdRestorePhaseCompleted ||
		restore.Status.Phase == etcdv1alpha1.EtcdRestorePhaseFailed {
		// Do nothing here. We're already finished.
		return ctrl.Result{}, nil
	} else {
		// In any other state continue with the reconciliation. In particular we don't *read* the other possible states
		// as once we know that we need to reconcile at all we will use the observed sate of the cluster and not our own
		// status field.
	}

	// Simulate the cluster, peer, and PVC we'll create. We won't ever *create* any of these other than the PVC, but
	// we need them to generate the PVC using the same code-path that the main operator does. Yes, in theory if you
	// upgraded the operator part-way through a restore things could get super fun if the way of making these changed.
	expectedCluster := clusterForRestore(restore)
	// We know that the cluster will be of size one for the restore. So only consider peer zero.
	expectedPeerName := expectedPeerNamesForCluster(expectedCluster)[0]
	expectedPeer := peerForCluster(expectedCluster, expectedPeerName)
	expectedPVC := pvcForPeer(expectedPeer)
	// We want to avoid a risk of conflicting either with an existing cluster or another restore targeted at the same
	// cluster name. So we'll mark the PVC's metadata to know that it is ours.
	markPVC(restore, expectedPVC)

	// Step 1: If the PVC does not exist, create it.
	var pvc corev1.PersistentVolumeClaim
	err := r.Get(ctx, name(expectedPVC), &pvc)

	if apierrors.IsNotFound(err) {
		// No PVC. Make it!
		err := r.Create(ctx, expectedPVC)
		return ctrl.Result{}, err
	} else if err != nil {
		// There was some other, non non-found error. Exit as we can't handle this case.
		return ctrl.Result{}, err
	} else {
		// The PVC is there already. Continue.
	}

	// Check to make sure we're expecting to restore into this PVC
	if !IsOurPVC(restore, pvc) {
		// It's unsafe to take any further action. We don't control the PVC.
		// TODO Report error reason
		restore.Status.Phase = etcdv1alpha1.EtcdRestorePhaseFailed
		err := r.Client.Status().Update(ctx, &restore)
		return ctrl.Result{}, err
	} else {
		// Great, this PVC is marked as *our* restore (and not any random PVC, or one for an existing cluster, or
		// a conflicting restore).
	}

	// Step 2: If the PVC is *not* already bound and we can't see a successful restore job, restore into it.
	if pvc.Status.Phase == corev1.ClaimBound {
		// PVC is already bound to something. This could be an already-running etcd Pod, or a restore Pod that's still
		// running. So leave it alone and come back later. We assume that this is probably a restore Pod that we started
		// the last time, so don't do anything and we will re-reconcile when the Pod exists anyway.
		return ctrl.Result{}, err
	} else if pvc.Status.Phase == corev1.ClaimLost {
		// Huh. If the claim is lost then the underlying PV is gone. This is beyond us, so just give up and let the user
		// figure it out.
		// TODO Report error reason
		restore.Status.Phase = etcdv1alpha1.EtcdRestorePhaseFailed
		err := r.Client.Status().Update(ctx, &restore)
		return ctrl.Result{}, err
	} else if pvc.Status.Phase == corev1.ClaimPending {
		// Okay, the PVC exists but is unbound. Continue.
	}

	// Check to see if our restore Pod already exists. It'll have the same name as the `EtcdRestore`.
	restorePod := corev1.Pod{}
	err = r.Get(ctx, name(&restore), &restorePod)
	if apierrors.IsNotFound(err) {
		// No Pod. Launch it!
		err := r.Create(ctx, podForRestore(restore, pvc.Name))
		return ctrl.Result{}, err
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Look for our labels if those aren't there this could be a pre-existing Pod and everything is bad.
	if !IsOurPod(restore, restorePod) {
		// TODO Report error reason
		restore.Status.Phase = etcdv1alpha1.EtcdRestorePhaseFailed
		err := r.Client.Status().Update(ctx, &restore)
		return ctrl.Result{}, err
	}

	phase := restorePod.Status.Phase
	if phase == corev1.PodSucceeded {
		// Success is what we want. Continue with reconciliation.
	} else if phase == corev1.PodFailed || phase == corev1.PodUnknown {
		// Pod has failed, so we fail
		// TODO Report error reason
		restore.Status.Phase = etcdv1alpha1.EtcdRestorePhaseFailed
		err := r.Client.Status().Update(ctx, &restore)
		return ctrl.Result{}, err
	} else {
		// This covers the "Pending" and "Running" phases. Do nothing and wait for the Pod to finish.
		return ctrl.Result{}, nil
	}

	// We've done the restore, the pod has exited. We can now create a single node cluster to sit atop the PVC we
	// already created.
	cluster := etcdv1alpha1.EtcdCluster{}
	err = r.Get(ctx, name(expectedCluster), &cluster)
	if apierrors.IsNotFound(err) {
		// No Cluster. Create it
		err := r.Create(ctx, expectedCluster)
		return ctrl.Result{}, err
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// If the cluster we just found isn't restored by us then fail completely
	if !IsOurCluster(restore, &cluster) {
		// TODO Report error reason
		restore.Status.Phase = etcdv1alpha1.EtcdRestorePhaseFailed
		err := r.Client.Status().Update(ctx, &restore)
		return ctrl.Result{}, err
	}

	// TODO Scale cluster to desired size. Need to check that the cluster is stable *before* triggering a scale. Then
	// wait for it to stabilise again before calling it a good job and exiting successfully.

	return ctrl.Result{}, nil
}

func podForRestore(restore etcdv1alpha1.EtcdRestore, pvcName string) *corev1.Pod {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      restore.Name,
			Namespace: restore.Namespace,
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "etcd-data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
						},
					},
				},
				{
					// This is scratch space to download the snapshot file into
					Name: "snapshot",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: nil,
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name: restoreContainerName,
					// Use the same etcd image that we actually run etcd with.
					Image: etcdRestoreImage,
					Env: []corev1.EnvVar{
						{
							Name:  "RESTORE_ETCD_DATA_DIR",
							Value: "/var/lib/etcd",
						},
						{
							Name:  "RESTORE_SNAPSHOT_DIR",
							Value: "/var/snapshot",
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "etcd-data",
							ReadOnly:  false,
							MountPath: "/var/lib/etcd",
						},
						{
							Name:      "snapshot",
							ReadOnly:  false,
							MountPath: "/var/snapshot",
						},
					},
				},
			},
			// If we fail, we fail
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	// Create a helper to append environment variables to the restore container.
	addEnvVar := func(ev corev1.EnvVar) {
		// We know there's only one container and it's the first in the list
		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, ev)
	}

	if restore.Spec.Source.GCSBucket != nil {
		// GCS restore
		addEnvVar(corev1.EnvVar{
			Name:  "RESTORE_TYPE",
			Value: "GCS"})
		addEnvVar(corev1.EnvVar{
			Name:  "RESTORE_GCS_BUCKET_NAME",
			Value: restore.Spec.Source.GCSBucket.BucketName})
		addEnvVar(corev1.EnvVar{
			Name:  "RESTORE_GCS_OBJECT_PATH",
			Value: restore.Spec.Source.GCSBucket.ObjectPath})
		addEnvVar(corev1.EnvVar{
			Name: "RESTORE_GCS_SECRET_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: restore.Spec.Source.GCSBucket.Credentials.SecretKeyRef}})
	}

	markPod(restore, &pod)
	return &pod
}

func clusterForRestore(restore etcdv1alpha1.EtcdRestore) *etcdv1alpha1.EtcdCluster {
	cluster := &etcdv1alpha1.EtcdCluster{
		ObjectMeta: v1.ObjectMeta{
			Name:      restore.Spec.ClusterTemplate.ClusterName,
			Namespace: restore.Namespace,
		},
		Spec: restore.Spec.ClusterTemplate.Spec,
	}
	// Always override to one. We'll make it bigger later if we need to.
	cluster.Spec.Replicas = pointer.Int32Ptr(1)
	// Slap a label on it so we can see it later
	markCluster(restore, cluster)
	return cluster
}

func (r *EtcdRestoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&etcdv1alpha1.EtcdRestore{}).
		Complete(r)
}
