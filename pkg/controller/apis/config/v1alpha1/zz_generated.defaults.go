// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by defaulter-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	v1alpha1 "k8s.io/kube-controller-manager/config/v1alpha1"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&v1alpha1.KubeControllerManagerConfiguration{}, func(obj interface{}) {
		SetObjectDefaults_KubeControllerManagerConfiguration(obj.(*v1alpha1.KubeControllerManagerConfiguration))
	})
	return nil
}

func SetObjectDefaults_KubeControllerManagerConfiguration(in *v1alpha1.KubeControllerManagerConfiguration) {
	SetDefaults_KubeControllerManagerConfiguration(in)
	SetDefaults_KubeCloudSharedConfiguration(&in.KubeCloudShared)
	SetDefaults_VolumeConfiguration(&in.PersistentVolumeBinderController.VolumeConfiguration)
	SetDefaults_PersistentVolumeRecyclerConfiguration(&in.PersistentVolumeBinderController.VolumeConfiguration.PersistentVolumeRecyclerConfiguration)
}
