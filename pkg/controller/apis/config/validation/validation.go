/*
Copyright 2018 The Kubernetes Authors.

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

package validation

import (
	"fmt"
	"net"
	"strings"

	apimachinery "k8s.io/apimachinery/pkg/apis/config/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	apiserver "k8s.io/apiserver/pkg/apis/config/validation"
	"k8s.io/kubernetes/pkg/controller/apis/config"
)

// ValidateKubeControllerManagerConfiguration ensures validation of the ValidateKubeControllerManagerConfiguration struct
func ValidateKubeControllerManagerConfiguration(obj *config.KubeControllerManagerConfiguration, allControllers, disabledByDefaultControllers []string) field.ErrorList {
	allErrs := field.ErrorList{}
	newPath := field.NewPath("KubeControllerManagerConfiguration")

	// add validate for GenericControllerManagerConfiguration values
	allErrs = append(allErrs, ValidateGenericControllerManagerConfiguration(&obj.Generic, allControllers, disabledByDefaultControllers, newPath.Child("generic"))...)
	// add validate for KubeCloudSharedConfiguration values
	allErrs = append(allErrs, ValidateKubeCloudSharedConfiguration(&obj.KubeCloudShared, newPath.Child("kubeCloudShared"))...)
	// add validate for ServiceControllerConfiguration values
	allErrs = append(allErrs, ValidateServiceControllerConfiguration(&obj.ServiceController, newPath.Child("serviceController"))...)
	// add validate for AttachDetachControllerConfiguration values
	allErrs = append(allErrs, ValidateAttachDetachControllerConfiguration(&obj.AttachDetachController, newPath.Child("attachDetachController"))...)
	// add validate for CSRSigningControllerConfiguration values
	allErrs = append(allErrs, ValidateCSRSigningControllerConfiguration(&obj.CSRSigningController, newPath.Child("csrSigningController"))...)
	// add validate for DeploymentControllerConfiguration values
	allErrs = append(allErrs, ValidateDeploymentControllerConfiguration(&obj.DeploymentController, newPath.Child("deploymentController"))...)
	// add validate for DaemonSetControllerConfiguration values
	allErrs = append(allErrs, ValidateDaemonSetControllerConfiguration(&obj.DaemonSetController, newPath.Child("daemonSetController"))...)
	// add validate for DeprecatedControllerConfiguration values
	allErrs = append(allErrs, ValidateDeprecatedControllerConfiguration(&obj.DeprecatedController, newPath.Child("deprecatedController"))...)
	// add validate for EndpointControllerConfiguration values
	allErrs = append(allErrs, ValidateEndpointControllerConfiguration(&obj.EndpointController, newPath.Child("endpointController"))...)
	// add validate for GarbageCollectorControllerConfiguration values
	allErrs = append(allErrs, ValidateGarbageCollectorControllerConfiguration(&obj.GarbageCollectorController, newPath.Child("garbageCollectorController"))...)
	// add validate for HPAControllerConfiguration values
	allErrs = append(allErrs, ValidateHPAControllerConfiguration(&obj.HPAController, newPath.Child("hpaController"))...)
	// add validate for JobControllerConfiguration values
	allErrs = append(allErrs, ValidateJobControllerConfiguration(&obj.JobController, newPath.Child("jobController"))...)
	// add validate for NamespaceControllerConfiguration values
	allErrs = append(allErrs, ValidateNamespaceControllerConfiguration(&obj.NamespaceController, newPath.Child("namespaceController"))...)
	// add validate for NodeIPAMControllerConfiguration values
	allErrs = append(allErrs, ValidateNodeIPAMControllerConfiguration(&obj.NodeIPAMController, newPath.Child("nodeIPAMController"))...)
	// add validate for NodeLifecycleControllerConfiguration values
	allErrs = append(allErrs, ValidateNodeLifecycleControllerConfiguration(&obj.NodeLifecycleController, newPath.Child("nodeLifecycleController"))...)
	// add validate for PersistentVolumeBinderControllerConfiguration values
	allErrs = append(allErrs, ValidatePersistentVolumeBinderControllerConfiguration(&obj.PersistentVolumeBinderController, newPath.Child("persistentVolumeBinderController"))...)
	// add validate for PodGCControllerConfiguration values
	allErrs = append(allErrs, ValidatePodGCControllerConfiguration(&obj.PodGCController, newPath.Child("podGCController"))...)
	// add validate for ReplicaSetControllerConfiguration values
	allErrs = append(allErrs, ValidateReplicaSetControllerConfiguration(&obj.ReplicaSetController, newPath.Child("replicaSetController"))...)
	// add validate for ReplicationControllerConfiguration values
	allErrs = append(allErrs, ValidateReplicationControllerConfiguration(&obj.ReplicationController, newPath.Child("replicationController"))...)
	// add validate for ResourceQuotaControllerConfiguration values
	allErrs = append(allErrs, ValidateResourceQuotaControllerConfiguration(&obj.ResourceQuotaController, newPath.Child("resourceQuotaController"))...)
	// add validate for SAControllerConfiguration values
	allErrs = append(allErrs, ValidateSAControllerConfigurationration(&obj.SAController, newPath.Child("saController"))...)
	// add validate for TTLAfterFinishedControllerConfiguration values
	allErrs = append(allErrs, ValidateTTLAfterFinishedControllerConfiguration(&obj.TTLAfterFinishedController, newPath.Child("ttlAfterFinishedController"))...)

	return allErrs
}

// ValidateGenericControllerManagerConfiguration ensures validation of the GenericControllerManagerConfiguration struct
func ValidateGenericControllerManagerConfiguration(obj *config.GenericControllerManagerConfiguration, allControllers, disabledByDefaultControllers []string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if ip := net.ParseIP(obj.Address); ip == nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("address"), obj.Address, "must be a valid IP"))
	}
	if obj.MinResyncPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("minResyncPeriod"), obj.MinResyncPeriod, "must be greater than zero"))
	}
	allErrs = append(allErrs, apimachinery.ValidateClientConnectionConfiguration(&obj.ClientConnection, fldPath.Child("clientConnection"))...)
	if obj.ControllerStartInterval.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("controllerStartInterval"), obj.ControllerStartInterval, "must be greater than zero"))
	}
	if !obj.LeaderElection.LeaderElect {
		return allErrs
	}
	allErrs = append(allErrs, apiserver.ValidateLeaderElectionConfiguration(&obj.LeaderElection, fldPath.Child("leaderElectionConfiguration"))...)
	allControllersSet := sets.NewString(allControllers...)
	for _, controller := range obj.Controllers {
		if controller == "*" {
			continue
		}
		if strings.HasPrefix(controller, "-") {
			controller = controller[1:]
		}
		if !allControllersSet.Has(controller) {
			errMsg := fmt.Sprintf("%q is not in the list of known controllers", controller)
			allErrs = append(allErrs, field.Invalid(fldPath.Child("controller"), obj.Controllers, errMsg))
		}
	}

	return allErrs
}

// ValidateKubeCloudSharedConfiguration ensures validation of the KubeCloudSharedConfiguration struct
func ValidateKubeCloudSharedConfiguration(obj *config.KubeCloudSharedConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if !obj.ConfigureCloudRoutes {
		return allErrs
	}
	if len(obj.ClusterName) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("clusterName"), ""))
	}
	if obj.NodeMonitorPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("nodeMonitorPeriod"), obj.NodeMonitorPeriod, "must be greater than zero"))
	}
	if obj.RouteReconciliationPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("routeReconciliationPeriod"), obj.RouteReconciliationPeriod, "must be greater than zero"))
	}

	return allErrs
}

// ValidateServiceControllerConfiguration ensures validation of the ServiceControllerConfiguration struct
func ValidateServiceControllerConfiguration(obj *config.ServiceControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentServiceSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentServiceSyncs"), obj.ConcurrentServiceSyncs, "must be greater than zero"))
	}

	return allErrs
}

// ValidateAttachDetachControllerConfiguration ensures validation of the AttachDetachControllerConfiguration struct
func ValidateAttachDetachControllerConfiguration(obj *config.AttachDetachControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ReconcilerSyncLoopPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("reconcilerSyncLoopPeriod"), obj.ReconcilerSyncLoopPeriod, "must be greater than zero"))
	}

	return allErrs
}

// ValidateCSRSigningControllerConfiguration ensures validation of the CSRSigningControllerConfiguration struct
func ValidateCSRSigningControllerConfiguration(obj *config.CSRSigningControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ClusterSigningDuration.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterSigningDuration"), obj.ClusterSigningDuration, "must be greater than zero"))
	}
	if len(obj.ClusterSigningCertFile) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("clusterSigningCertFile"), ""))
	}
	if len(obj.ClusterSigningKeyFile) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("clusterSigningKeyFile"), ""))
	}

	return allErrs
}

// ValidateDeploymentControllerConfiguration ensures validation of the DeploymentControllerConfiguration struct
func ValidateDeploymentControllerConfiguration(obj *config.DeploymentControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if obj.ConcurrentDeploymentSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentDeploymentSyncs"), obj.ConcurrentDeploymentSyncs, "must be greater than zero"))
	}

	if obj.DeploymentControllerSyncPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("deploymentControllerSyncPeriod"), obj.DeploymentControllerSyncPeriod, "must be greater than zero"))
	}

	return allErrs
}

// ValidateDaemonSetControllerConfiguration ensures validation of the DaemonSetControllerConfiguration struct
func ValidateDaemonSetControllerConfiguration(obj *config.DaemonSetControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentDaemonSetSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentDaemonSetSyncs"), obj.ConcurrentDaemonSetSyncs, "must be greater than zero"))
	}

	return allErrs
}

// ValidateDeprecatedControllerConfiguration ensures validation of the DeprecatedControllerConfiguration struct
func ValidateDeprecatedControllerConfiguration(obj *config.DeprecatedControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.RegisterRetryCount <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("registerRetryCount"), obj.RegisterRetryCount, "must be greater than zero"))
	}

	return allErrs
}

// ValidateEndpointControllerConfiguration ensures validation of the EndpointControllerConfiguration struct
func ValidateEndpointControllerConfiguration(obj *config.EndpointControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentEndpointSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentEndpointSyncs"), obj.ConcurrentEndpointSyncs, "must be greater than zero"))
	}

	return allErrs
}

// ValidateGarbageCollectorControllerConfiguration ensures validation of the GarbageCollectorControllerConfiguration struct
func ValidateGarbageCollectorControllerConfiguration(obj *config.GarbageCollectorControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if !obj.EnableGarbageCollector {
		return allErrs
	}
	if obj.ConcurrentGCSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentGCSyncs"), obj.ConcurrentGCSyncs, "must be greater than zero"))
	}

	return allErrs
}

// ValidateHPAControllerConfiguration ensures validation of the HPAControllerConfiguration struct
func ValidateHPAControllerConfiguration(obj *config.HPAControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if !obj.HorizontalPodAutoscalerUseRESTClients {
		return allErrs
	}
	if obj.HorizontalPodAutoscalerSyncPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("horizontalPodAutoscalerSyncPeriod"), obj.HorizontalPodAutoscalerSyncPeriod, "must be greater than zero"))
	}
	if obj.HorizontalPodAutoscalerUpscaleForbiddenWindow.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("horizontalPodAutoscalerUpscaleForbiddenWindow"), obj.HorizontalPodAutoscalerUpscaleForbiddenWindow, "must be greater than zero"))
	}
	if obj.HorizontalPodAutoscalerDownscaleStabilizationWindow.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("HorizontalPodAutoscalerDownscaleStabilizationWindow"), obj.HorizontalPodAutoscalerDownscaleStabilizationWindow, "must be greater than zero"))
	}
	if obj.HorizontalPodAutoscalerCPUInitializationPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("horizontalPodAutoscalerCPUInitializationPeriod"), obj.HorizontalPodAutoscalerCPUInitializationPeriod, "must be greater than zero"))
	}
	if obj.HorizontalPodAutoscalerInitialReadinessDelay.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("horizontalPodAutoscalerInitialReadinessDelay"), obj.HorizontalPodAutoscalerInitialReadinessDelay, "must be greater than zero"))
	}
	if obj.HorizontalPodAutoscalerDownscaleForbiddenWindow.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("horizontalPodAutoscalerDownscaleForbiddenWindow"), obj.HorizontalPodAutoscalerDownscaleForbiddenWindow, "must be greater than zero"))
	}
	if obj.HorizontalPodAutoscalerTolerance <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("horizontalPodAutoscalerTolerance"), obj.HorizontalPodAutoscalerTolerance, "must be greater than zero"))
	}

	return allErrs
}

// ValidateJobControllerConfiguration ensures validation of the JobControllerConfiguration struct
func ValidateJobControllerConfiguration(obj *config.JobControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentJobSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentJobSyncs"), obj.ConcurrentJobSyncs, "must be greater than zero"))
	}

	return allErrs
}

// ValidateNamespaceControllerConfiguration ensures validation of the NamespaceControllerConfiguration struct
func ValidateNamespaceControllerConfiguration(obj *config.NamespaceControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentNamespaceSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentNamespaceSyncs"), obj.ConcurrentNamespaceSyncs, "must be greater than zero"))
	}
	if obj.NamespaceSyncPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("namespaceSyncPeriod"), obj.NamespaceSyncPeriod, "must be greater than zero"))
	}

	return allErrs
}

// ValidateNodeIPAMControllerConfiguration ensures validation of the NodeIPAMControllerConfiguration struct
func ValidateNodeIPAMControllerConfiguration(obj *config.NodeIPAMControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.NodeCIDRMaskSize <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("nodeCIDRMaskSize"), obj.NodeCIDRMaskSize, "must be greater than zero"))
	}

	return allErrs
}

// ValidateNodeLifecycleControllerConfiguration ensures validation of the NodeLifecycleControllerConfiguration struct
func ValidateNodeLifecycleControllerConfiguration(obj *config.NodeLifecycleControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if !obj.EnableTaintManager {
		return allErrs
	}
	if obj.PodEvictionTimeout.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("podEvictionTimeout"), obj.PodEvictionTimeout, "must be greater than zero"))
	}
	if obj.NodeMonitorGracePeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("nodeMonitorGracePeriod"), obj.NodeMonitorGracePeriod, "must be greater than zero"))
	}
	if obj.NodeStartupGracePeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("nodeStartupGracePeriod"), obj.NodeStartupGracePeriod, "must be greater than zero"))
	}

	return allErrs
}

// ValidatePersistentVolumeBinderControllerConfiguration ensures validation of the PersistentVolumeBinderControllerConfiguration struct
func ValidatePersistentVolumeBinderControllerConfiguration(obj *config.PersistentVolumeBinderControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.PVClaimBinderSyncPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("pvclaimBinderSyncPeriod"), obj.PVClaimBinderSyncPeriod, "must be greater than zero"))
	}
	allErrs = append(allErrs, ValidateVolumeConfiguration(&obj.VolumeConfiguration, fldPath.Child("volumeConfiguration"))...)

	return allErrs
}

// ValidatePodGCControllerConfiguration ensures validation of the PodGCControllerConfiguration struct
func ValidatePodGCControllerConfiguration(obj *config.PodGCControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.TerminatedPodGCThreshold <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("terminatedPodGCThreshold"), obj.TerminatedPodGCThreshold, "must be greater than zero"))
	}

	return allErrs
}

// ValidateReplicaSetControllerConfiguration ensures validation of the ReplicaSetControllerConfiguration struct
func ValidateReplicaSetControllerConfiguration(obj *config.ReplicaSetControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentRSSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentRSSyncs"), obj.ConcurrentRSSyncs, "must be greater than zero"))
	}

	return allErrs
}

// ValidateReplicationControllerConfiguration ensures validation of the ReplicationControllerConfiguration struct
func ValidateReplicationControllerConfiguration(obj *config.ReplicationControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentRCSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentRCSyncs"), obj.ConcurrentRCSyncs, "must be greater than zero"))
	}

	return allErrs
}

// ValidateResourceQuotaControllerConfiguration ensures validation of the ResourceQuotaControllerConfiguration struct
func ValidateResourceQuotaControllerConfiguration(obj *config.ResourceQuotaControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentResourceQuotaSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentResourceQuotaSyncs"), obj.ConcurrentResourceQuotaSyncs, "must be greater than zero"))
	}
	if obj.ResourceQuotaSyncPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("resourceQuotaSyncPeriod"), obj.ResourceQuotaSyncPeriod, "must be greater than zero"))
	}

	return allErrs
}

// ValidateSAControllerConfigurationration ensures validation of the SAControllerConfiguration struct
func ValidateSAControllerConfigurationration(obj *config.SAControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentSATokenSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentSATokenSyncs"), obj.ConcurrentSATokenSyncs, "must be greater than zero"))
	}

	return allErrs
}

// ValidateTTLAfterFinishedControllerConfiguration ensures validation of the TTLAfterFinishedControllerConfiguration struct
func ValidateTTLAfterFinishedControllerConfiguration(obj *config.TTLAfterFinishedControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.ConcurrentTTLSyncs <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("concurrentTTLSyncs"), obj.ConcurrentTTLSyncs, "must be greater than zero"))
	}

	return allErrs
}

// ValidateVolumeConfiguration ensures validation of the VolumeConfiguration struct
func ValidateVolumeConfiguration(obj *config.VolumeConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if !obj.EnableDynamicProvisioning || obj.EnableHostPathProvisioning {
		return allErrs
	}
	if len(obj.FlexVolumePluginDir) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("flexVolumePluginDir"), ""))
	}
	allErrs = append(allErrs, ValidatePersistentVolumeRecyclerConfiguration(&obj.PersistentVolumeRecyclerConfiguration, fldPath.Child("persistentVolumeRecyclerConfiguration"))...)

	return allErrs
}

// ValidatePersistentVolumeRecyclerConfiguration ensures validation of the PersistentVolumeRecyclerConfiguration struct
func ValidatePersistentVolumeRecyclerConfiguration(obj *config.PersistentVolumeRecyclerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if obj.MaximumRetry <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("maximumRetry"), obj.MaximumRetry, "must be greater than zero"))
	}
	if obj.MinimumTimeoutNFS <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("minimumTimeoutNFS"), obj.MinimumTimeoutNFS, "must be greater than zero"))
	}
	if obj.IncrementTimeoutNFS <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("incrementTimeoutNFS"), obj.IncrementTimeoutNFS, "must be greater than zero"))
	}
	if obj.MinimumTimeoutHostPath <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("minimumTimeoutHostPath"), obj.MinimumTimeoutHostPath, "must be greater than zero"))
	}
	if obj.IncrementTimeoutHostPath <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("incrementTimeoutHostPath"), obj.IncrementTimeoutHostPath, "must be greater than zero"))
	}

	return allErrs
}
