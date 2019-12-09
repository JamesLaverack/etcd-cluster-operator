package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type EtcdRestoreMethod struct {
	// Restore from an etcd snapshot
	// +optional
	Snapshot *EtcdRestoreSnapshotLocation `json:"snapshot,omitempty"`
}

type EtcdRestoreSnapshotLocation struct {
	// GCSBucket identifies a Google Cloud Storage bucket to pull the snapshot from
	// +optional
	GCSBucket *EtcdRestoreGCSBucket `json:"gcsBucket,omitempty"`
}

type EtcdRestoreGCSBucket struct {
	// BucketName is the name of the storage bucket.
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=222
	BucketName string `json:"bucketName"`

	// ObjectPath is the path to the object inside the bucket.
	// +kubebuilder:validation:MinLength=1
	ObjectPath string `json:"objectPath"`

	// Credentials holds the method of obtaining credentials that will be provided to the
	// Google Cloud APIs in order to write backup data.
	// +optional
	Credentials *GoogleCloudCredentials `json:"credentials,omitempty"`
}

// A template to define the cluster we'll make. The namespace will be the same as this restore resource.
type EtcdClusterTemplate struct {
	// ClusterName is the name of the EtcdCluster that will be created by this operation.
	ClusterName string `json:"clusterName"`

	// Spec is the specification of the cluster that will be created by this operation
	Spec EtcdClusterSpec `json:"spec"`
}

// EtcdRestoreSpec defines the desired state of EtcdRestore
type EtcdRestoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Restore EtcdRestoreMethod `json:"restore"`

	ClusterTemplate EtcdClusterTemplate `json:"clusterTemplate"`
}

type EtcdRestorePhase string

var (
	EtcdRestorePhaseCreatingPVCs      EtcdRestorePhase = "CreatingPVCs"
	EtcdRestorePhaseUnpackingSnapshot EtcdRestorePhase = "UnpackingSnapshot"
	EtcdRestorePhaseCreatingCluster   EtcdRestorePhase = "CreatingCluster"
	EtcdRestorePhaseFailed            EtcdRestorePhase = "Failed"
	EtcdRestorePhaseCompleted         EtcdRestorePhase = "Completed"
)

// EtcdRestoreStatus defines the observed state of EtcdRestore
type EtcdRestoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase is what the restore is doing last time it was checked. The possible end states are "Failed" and
	// "Completed".
	Phase EtcdRestorePhase `json:"phase,omitempty"`
}

// +kubebuilder:object:root=true

// EtcdRestore is the Schema for the etcdrestores API
type EtcdRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdRestoreSpec   `json:"spec,omitempty"`
	Status EtcdRestoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EtcdRestoreList contains a list of EtcdRestore
type EtcdRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdRestore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdRestore{}, &EtcdRestoreList{})
}
