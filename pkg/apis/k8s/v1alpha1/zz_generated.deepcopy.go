// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatastoreSpec) DeepCopyInto(out *DatastoreSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatastoreSpec.
func (in *DatastoreSpec) DeepCopy() *DatastoreSpec {
	if in == nil {
		return nil
	}
	out := new(DatastoreSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Grackle) DeepCopyInto(out *Grackle) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Grackle.
func (in *Grackle) DeepCopy() *Grackle {
	if in == nil {
		return nil
	}
	out := new(Grackle)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Grackle) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrackleList) DeepCopyInto(out *GrackleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Grackle, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrackleList.
func (in *GrackleList) DeepCopy() *GrackleList {
	if in == nil {
		return nil
	}
	out := new(GrackleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GrackleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrackleSpec) DeepCopyInto(out *GrackleSpec) {
	*out = *in
	out.Datastore = in.Datastore
	if in.Ingest != nil {
		in, out := &in.Ingest, &out.Ingest
		*out = new(IngestSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Web != nil {
		in, out := &in.Web, &out.Web
		*out = new(WebSpec)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrackleSpec.
func (in *GrackleSpec) DeepCopy() *GrackleSpec {
	if in == nil {
		return nil
	}
	out := new(GrackleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrackleStatus) DeepCopyInto(out *GrackleStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrackleStatus.
func (in *GrackleStatus) DeepCopy() *GrackleStatus {
	if in == nil {
		return nil
	}
	out := new(GrackleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngestSpec) DeepCopyInto(out *IngestSpec) {
	*out = *in
	if in.Track != nil {
		in, out := &in.Track, &out.Track
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngestSpec.
func (in *IngestSpec) DeepCopy() *IngestSpec {
	if in == nil {
		return nil
	}
	out := new(IngestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebSpec) DeepCopyInto(out *WebSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebSpec.
func (in *WebSpec) DeepCopy() *WebSpec {
	if in == nil {
		return nil
	}
	out := new(WebSpec)
	in.DeepCopyInto(out)
	return out
}
