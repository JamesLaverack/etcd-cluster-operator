// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Bootstrap) DeepCopyInto(out *Bootstrap) {
	*out = *in
	if in.Static != nil {
		in, out := &in.Static, &out.Static
		*out = new(StaticBootstrap)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Bootstrap.
func (in *Bootstrap) DeepCopy() *Bootstrap {
	if in == nil {
		return nil
	}
	out := new(Bootstrap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackup) DeepCopyInto(out *EtcdBackup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackup.
func (in *EtcdBackup) DeepCopy() *EtcdBackup {
	if in == nil {
		return nil
	}
	out := new(EtcdBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdBackup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupDestination) DeepCopyInto(out *EtcdBackupDestination) {
	*out = *in
	if in.Local != nil {
		in, out := &in.Local, &out.Local
		*out = new(EtcdBackupDestinationLocal)
		**out = **in
	}
	if in.GCSBucket != nil {
		in, out := &in.GCSBucket, &out.GCSBucket
		*out = new(EtcdBackupDestinationGCSBucket)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupDestination.
func (in *EtcdBackupDestination) DeepCopy() *EtcdBackupDestination {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupDestination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupDestinationGCSBucket) DeepCopyInto(out *EtcdBackupDestinationGCSBucket) {
	*out = *in
	if in.Credentials != nil {
		in, out := &in.Credentials, &out.Credentials
		*out = new(GoogleCloudCredentials)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupDestinationGCSBucket.
func (in *EtcdBackupDestinationGCSBucket) DeepCopy() *EtcdBackupDestinationGCSBucket {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupDestinationGCSBucket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupDestinationLocal) DeepCopyInto(out *EtcdBackupDestinationLocal) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupDestinationLocal.
func (in *EtcdBackupDestinationLocal) DeepCopy() *EtcdBackupDestinationLocal {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupDestinationLocal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupList) DeepCopyInto(out *EtcdBackupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EtcdBackup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupList.
func (in *EtcdBackupList) DeepCopy() *EtcdBackupList {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdBackupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupSchedule) DeepCopyInto(out *EtcdBackupSchedule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupSchedule.
func (in *EtcdBackupSchedule) DeepCopy() *EtcdBackupSchedule {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupSchedule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdBackupSchedule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupScheduleList) DeepCopyInto(out *EtcdBackupScheduleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EtcdBackupSchedule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupScheduleList.
func (in *EtcdBackupScheduleList) DeepCopy() *EtcdBackupScheduleList {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupScheduleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdBackupScheduleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupScheduleSpec) DeepCopyInto(out *EtcdBackupScheduleSpec) {
	*out = *in
	in.BackupSpec.DeepCopyInto(&out.BackupSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupScheduleSpec.
func (in *EtcdBackupScheduleSpec) DeepCopy() *EtcdBackupScheduleSpec {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupScheduleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupScheduleStatus) DeepCopyInto(out *EtcdBackupScheduleStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupScheduleStatus.
func (in *EtcdBackupScheduleStatus) DeepCopy() *EtcdBackupScheduleStatus {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupScheduleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupSpec) DeepCopyInto(out *EtcdBackupSpec) {
	*out = *in
	if in.ClusterEndpoints != nil {
		in, out := &in.ClusterEndpoints, &out.ClusterEndpoints
		*out = make([]EtcdClusterEndpoint, len(*in))
		copy(*out, *in)
	}
	in.Destination.DeepCopyInto(&out.Destination)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupSpec.
func (in *EtcdBackupSpec) DeepCopy() *EtcdBackupSpec {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupStatus) DeepCopyInto(out *EtcdBackupStatus) {
	*out = *in
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = (*in).DeepCopy()
	}
	if in.CompletionTime != nil {
		in, out := &in.CompletionTime, &out.CompletionTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupStatus.
func (in *EtcdBackupStatus) DeepCopy() *EtcdBackupStatus {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdCluster) DeepCopyInto(out *EtcdCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdCluster.
func (in *EtcdCluster) DeepCopy() *EtcdCluster {
	if in == nil {
		return nil
	}
	out := new(EtcdCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdClusterEndpoint) DeepCopyInto(out *EtcdClusterEndpoint) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdClusterEndpoint.
func (in *EtcdClusterEndpoint) DeepCopy() *EtcdClusterEndpoint {
	if in == nil {
		return nil
	}
	out := new(EtcdClusterEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdClusterList) DeepCopyInto(out *EtcdClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EtcdCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdClusterList.
func (in *EtcdClusterList) DeepCopy() *EtcdClusterList {
	if in == nil {
		return nil
	}
	out := new(EtcdClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdClusterSpec) DeepCopyInto(out *EtcdClusterSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		*out = new(EtcdPeerStorage)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdClusterSpec.
func (in *EtcdClusterSpec) DeepCopy() *EtcdClusterSpec {
	if in == nil {
		return nil
	}
	out := new(EtcdClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdClusterStatus) DeepCopyInto(out *EtcdClusterStatus) {
	*out = *in
	if in.Members != nil {
		in, out := &in.Members, &out.Members
		*out = make([]EtcdMember, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdClusterStatus.
func (in *EtcdClusterStatus) DeepCopy() *EtcdClusterStatus {
	if in == nil {
		return nil
	}
	out := new(EtcdClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdClusterTemplate) DeepCopyInto(out *EtcdClusterTemplate) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdClusterTemplate.
func (in *EtcdClusterTemplate) DeepCopy() *EtcdClusterTemplate {
	if in == nil {
		return nil
	}
	out := new(EtcdClusterTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdMember) DeepCopyInto(out *EtcdMember) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdMember.
func (in *EtcdMember) DeepCopy() *EtcdMember {
	if in == nil {
		return nil
	}
	out := new(EtcdMember)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdPeer) DeepCopyInto(out *EtcdPeer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdPeer.
func (in *EtcdPeer) DeepCopy() *EtcdPeer {
	if in == nil {
		return nil
	}
	out := new(EtcdPeer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdPeer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdPeerList) DeepCopyInto(out *EtcdPeerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EtcdPeer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdPeerList.
func (in *EtcdPeerList) DeepCopy() *EtcdPeerList {
	if in == nil {
		return nil
	}
	out := new(EtcdPeerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdPeerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdPeerSpec) DeepCopyInto(out *EtcdPeerSpec) {
	*out = *in
	if in.Bootstrap != nil {
		in, out := &in.Bootstrap, &out.Bootstrap
		*out = new(Bootstrap)
		(*in).DeepCopyInto(*out)
	}
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		*out = new(EtcdPeerStorage)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdPeerSpec.
func (in *EtcdPeerSpec) DeepCopy() *EtcdPeerSpec {
	if in == nil {
		return nil
	}
	out := new(EtcdPeerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdPeerStatus) DeepCopyInto(out *EtcdPeerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdPeerStatus.
func (in *EtcdPeerStatus) DeepCopy() *EtcdPeerStatus {
	if in == nil {
		return nil
	}
	out := new(EtcdPeerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdPeerStorage) DeepCopyInto(out *EtcdPeerStorage) {
	*out = *in
	if in.VolumeClaimTemplate != nil {
		in, out := &in.VolumeClaimTemplate, &out.VolumeClaimTemplate
		*out = new(v1.PersistentVolumeClaimSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdPeerStorage.
func (in *EtcdPeerStorage) DeepCopy() *EtcdPeerStorage {
	if in == nil {
		return nil
	}
	out := new(EtcdPeerStorage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdRestore) DeepCopyInto(out *EtcdRestore) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdRestore.
func (in *EtcdRestore) DeepCopy() *EtcdRestore {
	if in == nil {
		return nil
	}
	out := new(EtcdRestore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdRestore) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdRestoreList) DeepCopyInto(out *EtcdRestoreList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EtcdRestore, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdRestoreList.
func (in *EtcdRestoreList) DeepCopy() *EtcdRestoreList {
	if in == nil {
		return nil
	}
	out := new(EtcdRestoreList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdRestoreList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdRestoreSource) DeepCopyInto(out *EtcdRestoreSource) {
	*out = *in
	if in.Local != nil {
		in, out := &in.Local, &out.Local
		*out = new(EtcdRestoreSourceLocal)
		**out = **in
	}
	if in.GCSBucket != nil {
		in, out := &in.GCSBucket, &out.GCSBucket
		*out = new(EtcdRestoreSourceGCSBucket)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdRestoreSource.
func (in *EtcdRestoreSource) DeepCopy() *EtcdRestoreSource {
	if in == nil {
		return nil
	}
	out := new(EtcdRestoreSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdRestoreSourceGCSBucket) DeepCopyInto(out *EtcdRestoreSourceGCSBucket) {
	*out = *in
	if in.Credentials != nil {
		in, out := &in.Credentials, &out.Credentials
		*out = new(GoogleCloudCredentials)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdRestoreSourceGCSBucket.
func (in *EtcdRestoreSourceGCSBucket) DeepCopy() *EtcdRestoreSourceGCSBucket {
	if in == nil {
		return nil
	}
	out := new(EtcdRestoreSourceGCSBucket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdRestoreSourceLocal) DeepCopyInto(out *EtcdRestoreSourceLocal) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdRestoreSourceLocal.
func (in *EtcdRestoreSourceLocal) DeepCopy() *EtcdRestoreSourceLocal {
	if in == nil {
		return nil
	}
	out := new(EtcdRestoreSourceLocal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdRestoreSpec) DeepCopyInto(out *EtcdRestoreSpec) {
	*out = *in
	in.Source.DeepCopyInto(&out.Source)
	in.ClusterTemplate.DeepCopyInto(&out.ClusterTemplate)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdRestoreSpec.
func (in *EtcdRestoreSpec) DeepCopy() *EtcdRestoreSpec {
	if in == nil {
		return nil
	}
	out := new(EtcdRestoreSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdRestoreStatus) DeepCopyInto(out *EtcdRestoreStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdRestoreStatus.
func (in *EtcdRestoreStatus) DeepCopy() *EtcdRestoreStatus {
	if in == nil {
		return nil
	}
	out := new(EtcdRestoreStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleCloudCredentials) DeepCopyInto(out *GoogleCloudCredentials) {
	*out = *in
	if in.SecretKeyRef != nil {
		in, out := &in.SecretKeyRef, &out.SecretKeyRef
		*out = new(v1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleCloudCredentials.
func (in *GoogleCloudCredentials) DeepCopy() *GoogleCloudCredentials {
	if in == nil {
		return nil
	}
	out := new(GoogleCloudCredentials)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitialClusterMember) DeepCopyInto(out *InitialClusterMember) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitialClusterMember.
func (in *InitialClusterMember) DeepCopy() *InitialClusterMember {
	if in == nil {
		return nil
	}
	out := new(InitialClusterMember)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StaticBootstrap) DeepCopyInto(out *StaticBootstrap) {
	*out = *in
	if in.InitialCluster != nil {
		in, out := &in.InitialCluster, &out.InitialCluster
		*out = make([]InitialClusterMember, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StaticBootstrap.
func (in *StaticBootstrap) DeepCopy() *StaticBootstrap {
	if in == nil {
		return nil
	}
	out := new(StaticBootstrap)
	in.DeepCopyInto(out)
	return out
}
