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
	"testing"
	"time"

	apimachineryconfig "k8s.io/apimachinery/pkg/apis/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	apiserverconfig "k8s.io/apiserver/pkg/apis/config"
	"k8s.io/kubernetes/pkg/controller/apis/config"
)

func TestValidateKubeControllerManagerConfiguration(t *testing.T) {
	validConfig := &config.KubeControllerManagerConfiguration{
		Generic: config.GenericControllerManagerConfiguration{
			Address:         "0.0.0.0",
			MinResyncPeriod: metav1.Duration{Duration: 8 * time.Hour},
			ClientConnection: apimachineryconfig.ClientConnectionConfiguration{
				Burst: 100,
			},
			ControllerStartInterval: metav1.Duration{Duration: 2 * time.Minute},
			LeaderElection: apiserverconfig.LeaderElectionConfiguration{
				ResourceLock:  "configmap",
				LeaderElect:   true,
				LeaseDuration: metav1.Duration{Duration: 30 * time.Second},
				RenewDeadline: metav1.Duration{Duration: 15 * time.Second},
				RetryPeriod:   metav1.Duration{Duration: 5 * time.Second},
			},
			Controllers: []string{"*"},
		},
		KubeCloudShared: config.KubeCloudSharedConfiguration{
			RouteReconciliationPeriod: metav1.Duration{Duration: 30 * time.Second},
			NodeMonitorPeriod:         metav1.Duration{Duration: 10 * time.Second},
			ConfigureCloudRoutes:      false,
		},
		ServiceController: config.ServiceControllerConfiguration{
			ConcurrentServiceSyncs: 2,
		},
		AttachDetachController: config.AttachDetachControllerConfiguration{
			ReconcilerSyncLoopPeriod: metav1.Duration{Duration: 30 * time.Second},
		},
		CSRSigningController: config.CSRSigningControllerConfiguration{
			ClusterSigningCertFile: "/cluster-signing-cert",
			ClusterSigningKeyFile:  "/cluster-signing-key",
			ClusterSigningDuration: metav1.Duration{Duration: 10 * time.Hour},
		},
		DaemonSetController: config.DaemonSetControllerConfiguration{
			ConcurrentDaemonSetSyncs: 2,
		},
		DeploymentController: config.DeploymentControllerConfiguration{
			ConcurrentDeploymentSyncs:      10,
			DeploymentControllerSyncPeriod: metav1.Duration{Duration: 45 * time.Second},
		},
		DeprecatedController: config.DeprecatedControllerConfiguration{
			RegisterRetryCount: 10,
		},
		EndpointController: config.EndpointControllerConfiguration{
			ConcurrentEndpointSyncs: 10,
		},
		GarbageCollectorController: config.GarbageCollectorControllerConfiguration{
			ConcurrentGCSyncs:      30,
			EnableGarbageCollector: true,
		},
		HPAController: config.HPAControllerConfiguration{
			HorizontalPodAutoscalerSyncPeriod:                   metav1.Duration{Duration: 45 * time.Second},
			HorizontalPodAutoscalerUpscaleForbiddenWindow:       metav1.Duration{Duration: 1 * time.Minute},
			HorizontalPodAutoscalerDownscaleForbiddenWindow:     metav1.Duration{Duration: 2 * time.Minute},
			HorizontalPodAutoscalerDownscaleStabilizationWindow: metav1.Duration{Duration: 3 * time.Minute},
			HorizontalPodAutoscalerCPUInitializationPeriod:      metav1.Duration{Duration: 90 * time.Second},
			HorizontalPodAutoscalerInitialReadinessDelay:        metav1.Duration{Duration: 50 * time.Second},
			HorizontalPodAutoscalerTolerance:                    0.1,
			HorizontalPodAutoscalerUseRESTClients:               true,
		},
		JobController: config.JobControllerConfiguration{
			ConcurrentJobSyncs: 5,
		},
		NamespaceController: config.NamespaceControllerConfiguration{
			NamespaceSyncPeriod:      metav1.Duration{Duration: 10 * time.Minute},
			ConcurrentNamespaceSyncs: 20,
		},
		NodeIPAMController: config.NodeIPAMControllerConfiguration{
			NodeCIDRMaskSize: 48,
		},
		NodeLifecycleController: config.NodeLifecycleControllerConfiguration{
			EnableTaintManager:     false,
			NodeMonitorGracePeriod: metav1.Duration{Duration: 30 * time.Second},
			NodeStartupGracePeriod: metav1.Duration{Duration: 30 * time.Second},
			PodEvictionTimeout:     metav1.Duration{Duration: 2 * time.Minute},
		},
		PersistentVolumeBinderController: config.PersistentVolumeBinderControllerConfiguration{
			PVClaimBinderSyncPeriod: metav1.Duration{Duration: 30 * time.Second},
			VolumeConfiguration: config.VolumeConfiguration{
				EnableDynamicProvisioning:  false,
				EnableHostPathProvisioning: true,
				FlexVolumePluginDir:        "/flex-volume-plugin",
				PersistentVolumeRecyclerConfiguration: config.PersistentVolumeRecyclerConfiguration{
					MaximumRetry:             3,
					MinimumTimeoutNFS:        200,
					IncrementTimeoutNFS:      45,
					MinimumTimeoutHostPath:   45,
					IncrementTimeoutHostPath: 45,
				},
			},
		},
		PodGCController: config.PodGCControllerConfiguration{
			TerminatedPodGCThreshold: 12000,
		},
		ReplicaSetController: config.ReplicaSetControllerConfiguration{
			ConcurrentRSSyncs: 10,
		},
		ReplicationController: config.ReplicationControllerConfiguration{
			ConcurrentRCSyncs: 10,
		},
		ResourceQuotaController: config.ResourceQuotaControllerConfiguration{
			ResourceQuotaSyncPeriod:      metav1.Duration{Duration: 10 * time.Minute},
			ConcurrentResourceQuotaSyncs: 10,
		},
		SAController: config.SAControllerConfiguration{
			ConcurrentSATokenSyncs: 10,
		},
		TTLAfterFinishedController: config.TTLAfterFinishedControllerConfiguration{
			ConcurrentTTLSyncs: 8,
		},
	}

	concurrentServiceSyncsLt0 := validConfig.DeepCopy()
	concurrentServiceSyncsLt0.ServiceController.ConcurrentServiceSyncs = -1

	reconcilerSyncLoopPeriodLt0 := validConfig.DeepCopy()
	reconcilerSyncLoopPeriodLt0.AttachDetachController.ReconcilerSyncLoopPeriod = metav1.Duration{Duration: -30 * time.Second}

	concurrentDaemonSetSyncsLt0 := validConfig.DeepCopy()
	concurrentDaemonSetSyncsLt0.DaemonSetController.ConcurrentDaemonSetSyncs = -2

	deploymentControllerInvalidate := validConfig.DeepCopy()
	deploymentControllerInvalidate.DeploymentController.ConcurrentDeploymentSyncs = -10
	deploymentControllerInvalidate.DeploymentController.DeploymentControllerSyncPeriod = metav1.Duration{Duration: -45 * time.Second}

	deprecatedControllerInvalidate := validConfig.DeepCopy()
	deprecatedControllerInvalidate.DeprecatedController.RegisterRetryCount = -10

	endpointControllerInvalidate := validConfig.DeepCopy()
	endpointControllerInvalidate.EndpointController.ConcurrentEndpointSyncs = -10

	garbageCollectorControllerInvalidate := validConfig.DeepCopy()
	garbageCollectorControllerInvalidate.GarbageCollectorController.ConcurrentGCSyncs = -10
	garbageCollectorControllerInvalidate.GarbageCollectorController.EnableGarbageCollector = true

	jobControllerInvalidate := validConfig.DeepCopy()
	jobControllerInvalidate.JobController.ConcurrentJobSyncs = -5

	namespaceControllerInvalidate := validConfig.DeepCopy()
	namespaceControllerInvalidate.NamespaceController.ConcurrentNamespaceSyncs = -10
	namespaceControllerInvalidate.NamespaceController.NamespaceSyncPeriod = metav1.Duration{Duration: -10 * time.Minute}

	nodeIPAMControllerInvalidate := validConfig.DeepCopy()
	nodeIPAMControllerInvalidate.NodeIPAMController.NodeCIDRMaskSize = -10

	podGCControllerInvalidate := validConfig.DeepCopy()
	podGCControllerInvalidate.PodGCController.TerminatedPodGCThreshold = -10

	replicaSetControllerInvalidate := validConfig.DeepCopy()
	replicaSetControllerInvalidate.ReplicaSetController.ConcurrentRSSyncs = -10

	replicationControllerInvalidate := validConfig.DeepCopy()
	replicationControllerInvalidate.ReplicationController.ConcurrentRCSyncs = -10

	resourceQuotaControllerInvalidate := validConfig.DeepCopy()
	resourceQuotaControllerInvalidate.ResourceQuotaController.ConcurrentResourceQuotaSyncs = -10
	resourceQuotaControllerInvalidate.ResourceQuotaController.ResourceQuotaSyncPeriod = metav1.Duration{Duration: -10 * time.Minute}

	saControllerInvalidate := validConfig.DeepCopy()
	saControllerInvalidate.SAController.ConcurrentSATokenSyncs = -10

	ttlAfterFinishedControllerInvalidate := validConfig.DeepCopy()
	ttlAfterFinishedControllerInvalidate.TTLAfterFinishedController.ConcurrentTTLSyncs = -8

	scenarios := map[string]struct {
		expectedToFail bool
		config         *config.KubeControllerManagerConfiguration
	}{
		"good": {
			expectedToFail: false,
			config:         validConfig,
		},
		"service-controller-configuration-invalid": {
			expectedToFail: true,
			config:         concurrentServiceSyncsLt0,
		},
		"attach-detach-controller-configuration-invalid": {
			expectedToFail: true,
			config:         reconcilerSyncLoopPeriodLt0,
		},
		"daemonset-controller-configuration-invalid": {
			expectedToFail: true,
			config:         concurrentDaemonSetSyncsLt0,
		},
		"deployment-configuration-invalid": {
			expectedToFail: true,
			config:         deploymentControllerInvalidate,
		},
		"deprecated-controller-configuration-invalid": {
			expectedToFail: true,
			config:         deprecatedControllerInvalidate,
		},
		"endpoint-controller-configuration-invalid": {
			expectedToFail: true,
			config:         endpointControllerInvalidate,
		},
		"garbage-controller-configuration-invalid": {
			expectedToFail: true,
			config:         garbageCollectorControllerInvalidate,
		},
		"job-controller-configuration-invalid": {
			expectedToFail: true,
			config:         jobControllerInvalidate,
		},
		"namespace-controller-configuration-invalid": {
			expectedToFail: true,
			config:         namespaceControllerInvalidate,
		},
		"nodeIPAM-controller-configuration-invalid": {
			expectedToFail: true,
			config:         nodeIPAMControllerInvalidate,
		},
		"pod-GC-controller-configuration-invalid": {
			expectedToFail: true,
			config:         podGCControllerInvalidate,
		},
		"replica-set-controller-configuration-invalid": {
			expectedToFail: true,
			config:         replicaSetControllerInvalidate,
		},
		"replication-controller-configuration-invalid": {
			expectedToFail: true,
			config:         replicationControllerInvalidate,
		},
		"resourcequota-controller-configuration-invalid": {
			expectedToFail: true,
			config:         resourceQuotaControllerInvalidate,
		},
		"sa-controller-configuration-invalid": {
			expectedToFail: true,
			config:         saControllerInvalidate,
		},
		"ttl-after-finished-controller-configuration-invalid": {
			expectedToFail: true,
			config:         ttlAfterFinishedControllerInvalidate,
		},
	}

	for name, scenario := range scenarios {
		errs := ValidateKubeControllerManagerConfiguration(scenario.config, []string{""}, []string{""})
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}
	}
}

func TestValidateGenericControllerManagerConfiguration(t *testing.T) {

	newPath := field.NewPath("KubeControllerManagerConfiguration")
	validConfig := &config.GenericControllerManagerConfiguration{
		Address:         "0.0.0.0",
		MinResyncPeriod: metav1.Duration{Duration: 8 * time.Hour},
		ClientConnection: apimachineryconfig.ClientConnectionConfiguration{
			Burst: 100,
		},
		ControllerStartInterval: metav1.Duration{Duration: 2 * time.Minute},
		LeaderElection: apiserverconfig.LeaderElectionConfiguration{
			ResourceLock:  "configmap",
			LeaderElect:   true,
			LeaseDuration: metav1.Duration{Duration: 30 * time.Second},
			RenewDeadline: metav1.Duration{Duration: 15 * time.Second},
			RetryPeriod:   metav1.Duration{Duration: 5 * time.Second},
		},
		Controllers: []string{"*"},
	}

	addrInvalid := validConfig.DeepCopy()
	addrInvalid.Address = "0.0.0.0.0"

	minResyncPeriodLt0 := validConfig.DeepCopy()
	minResyncPeriodLt0.MinResyncPeriod = metav1.Duration{Duration: -1 * time.Second}

	clientConnectionInvalidate := validConfig.DeepCopy()
	clientConnectionInvalidate.ClientConnection.Burst = -1

	controllerStartIntervalLt0 := validConfig.DeepCopy()
	controllerStartIntervalLt0.ControllerStartInterval = metav1.Duration{Duration: -1 * time.Second}

	scenarios := map[string]struct {
		expectedToFail bool
		config         *config.GenericControllerManagerConfiguration
	}{
		"good": {
			expectedToFail: false,
			config:         validConfig,
		},
		"address-invalid": {
			expectedToFail: true,
			config:         addrInvalid,
		},
		"min-resync-period-invalid": {
			expectedToFail: true,
			config:         minResyncPeriodLt0,
		},
		"client-connection-invalid": {
			expectedToFail: true,
			config:         clientConnectionInvalidate,
		},
		"controller-start-interval-invalid": {
			expectedToFail: true,
			config:         controllerStartIntervalLt0,
		},
	}

	for name, scenario := range scenarios {
		errs := ValidateGenericControllerManagerConfiguration(scenario.config, []string{""}, []string{""}, newPath.Child("generic"))
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}
	}
}

func TestValidateKubeCloudSharedConfiguration(t *testing.T) {

	newPath := field.NewPath("KubeControllerManagerConfiguration")
	validConfig := &config.KubeCloudSharedConfiguration{
		RouteReconciliationPeriod: metav1.Duration{Duration: 30 * time.Second},
		NodeMonitorPeriod:         metav1.Duration{Duration: 10 * time.Second},
		ClusterName:               "kubenetes",
		ConfigureCloudRoutes:      true,
	}

	routeReconciliationPeriodLt0 := validConfig.DeepCopy()
	routeReconciliationPeriodLt0.RouteReconciliationPeriod = metav1.Duration{Duration: -1 * time.Second}

	nodeMonitorPeriodLt0 := validConfig.DeepCopy()
	nodeMonitorPeriodLt0.RouteReconciliationPeriod = metav1.Duration{Duration: -1 * time.Second}

	clusterNameMissing := validConfig.DeepCopy()
	clusterNameMissing.ClusterName = ""

	scenarios := map[string]struct {
		expectedToFail bool
		config         *config.KubeCloudSharedConfiguration
	}{
		"good": {
			expectedToFail: false,
			config:         validConfig,
		},
		"route-reconciliation-period-lt-0": {
			expectedToFail: true,
			config:         routeReconciliationPeriodLt0,
		},
		"node-monitor-period-lt-0": {
			expectedToFail: true,
			config:         nodeMonitorPeriodLt0,
		},
		"cluster-name-missing": {
			expectedToFail: true,
			config:         clusterNameMissing,
		},
	}

	for name, scenario := range scenarios {
		errs := ValidateKubeCloudSharedConfiguration(scenario.config, newPath.Child("kubeCloudShared"))
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}

	}
}

func TestValidateCSRSigningControllerConfiguration(t *testing.T) {
	newPath := field.NewPath("KubeControllerManagerConfiguration")
	validConfig := &config.CSRSigningControllerConfiguration{
		ClusterSigningCertFile: "/cluster-signing-cert",
		ClusterSigningKeyFile:  "/cluster-signing-key",
		ClusterSigningDuration: metav1.Duration{Duration: 10 * time.Hour},
	}

	clusterSigningDurationLt0 := validConfig.DeepCopy()
	clusterSigningDurationLt0.ClusterSigningDuration = metav1.Duration{Duration: -1 * time.Second}

	clusterSigningCertFileNotExist := validConfig.DeepCopy()
	clusterSigningCertFileNotExist.ClusterSigningCertFile = ""

	clusterSigningKeyFileNotExist := validConfig.DeepCopy()
	clusterSigningKeyFileNotExist.ClusterSigningKeyFile = ""

	scenarios := map[string]struct {
		expectedToFail bool
		config         *config.CSRSigningControllerConfiguration
	}{
		"good": {
			expectedToFail: false,
			config:         validConfig,
		},
		"cluster-signingCertFile-not-exist": {
			expectedToFail: true,
			config:         clusterSigningCertFileNotExist,
		},
		"min-resync-period-invalid": {
			expectedToFail: true,
			config:         clusterSigningKeyFileNotExist,
		},
	}

	for name, scenario := range scenarios {
		errs := ValidateCSRSigningControllerConfiguration(scenario.config, newPath.Child("csrSigningController"))
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}

	}
}

func TestValidateHPAControllerConfiguration(t *testing.T) {
	newPath := field.NewPath("KubeControllerManagerConfiguration")
	validConfig := &config.HPAControllerConfiguration{
		HorizontalPodAutoscalerSyncPeriod:                   metav1.Duration{Duration: 45 * time.Second},
		HorizontalPodAutoscalerUpscaleForbiddenWindow:       metav1.Duration{Duration: 1 * time.Minute},
		HorizontalPodAutoscalerDownscaleForbiddenWindow:     metav1.Duration{Duration: 2 * time.Minute},
		HorizontalPodAutoscalerDownscaleStabilizationWindow: metav1.Duration{Duration: 3 * time.Minute},
		HorizontalPodAutoscalerCPUInitializationPeriod:      metav1.Duration{Duration: 90 * time.Second},
		HorizontalPodAutoscalerInitialReadinessDelay:        metav1.Duration{Duration: 50 * time.Second},
		HorizontalPodAutoscalerTolerance:                    0.1,
		HorizontalPodAutoscalerUseRESTClients:               true,
	}

	clusterSigningDurationLt0 := validConfig.DeepCopy()
	clusterSigningDurationLt0.HorizontalPodAutoscalerSyncPeriod = metav1.Duration{Duration: -1 * time.Second}

	HorizontalPodAutoscalerUpscaleForbiddenWindowLt0 := validConfig.DeepCopy()
	HorizontalPodAutoscalerUpscaleForbiddenWindowLt0.HorizontalPodAutoscalerUpscaleForbiddenWindow = metav1.Duration{Duration: -1 * time.Second}

	HorizontalPodAutoscalerDownscaleForbiddenWindowLt0 := validConfig.DeepCopy()
	HorizontalPodAutoscalerDownscaleForbiddenWindowLt0.HorizontalPodAutoscalerDownscaleForbiddenWindow = metav1.Duration{Duration: -1 * time.Second}

	HorizontalPodAutoscalerDownscaleStabilizationWindowLt0 := validConfig.DeepCopy()
	HorizontalPodAutoscalerDownscaleStabilizationWindowLt0.HorizontalPodAutoscalerDownscaleStabilizationWindow = metav1.Duration{Duration: -1 * time.Second}

	HorizontalPodAutoscalerCPUInitializationPeriodLt0 := validConfig.DeepCopy()
	HorizontalPodAutoscalerCPUInitializationPeriodLt0.HorizontalPodAutoscalerCPUInitializationPeriod = metav1.Duration{Duration: -1 * time.Second}

	HorizontalPodAutoscalerInitialReadinessDelayLt0 := validConfig.DeepCopy()
	HorizontalPodAutoscalerInitialReadinessDelayLt0.HorizontalPodAutoscalerInitialReadinessDelay = metav1.Duration{Duration: -1 * time.Second}

	scenarios := map[string]struct {
		expectedToFail bool
		config         *config.HPAControllerConfiguration
	}{
		"good": {
			expectedToFail: false,
			config:         validConfig,
		},
		"cluster-signing-duration-lt-0": {
			expectedToFail: true,
			config:         clusterSigningDurationLt0,
		},
		"horizontal-pod-autoscaler-upscale-forbidden-window-lt-0": {
			expectedToFail: true,
			config:         HorizontalPodAutoscalerUpscaleForbiddenWindowLt0,
		},
		"horizontal-pod-autoscaler-downscale-forbidden-window-lt-0": {
			expectedToFail: true,
			config:         HorizontalPodAutoscalerDownscaleForbiddenWindowLt0,
		},
		"horizontal-pod-autoscaler-downscale-stabilization-window-lt-0": {
			expectedToFail: true,
			config:         HorizontalPodAutoscalerDownscaleStabilizationWindowLt0,
		},
		"horizontal-pod-autoscaler-cpu-initialization-period-lt-0": {
			expectedToFail: true,
			config:         HorizontalPodAutoscalerCPUInitializationPeriodLt0,
		},
		"horizontal-pod-autoscaler-initial-readiness-delay-lt-0": {
			expectedToFail: true,
			config:         HorizontalPodAutoscalerInitialReadinessDelayLt0,
		},
	}

	for name, scenario := range scenarios {
		errs := ValidateHPAControllerConfiguration(scenario.config, newPath.Child("hpaController"))
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}
	}

}

func TestValidateNodeLifecycleControllerConfiguration(t *testing.T) {
	newPath := field.NewPath("KubeControllerManagerConfiguration")
	validConfig := &config.NodeLifecycleControllerConfiguration{
		EnableTaintManager:     true,
		NodeMonitorGracePeriod: metav1.Duration{Duration: 30 * time.Second},
		NodeStartupGracePeriod: metav1.Duration{Duration: 30 * time.Second},
		PodEvictionTimeout:     metav1.Duration{Duration: 2 * time.Minute},
	}

	nodeMonitorGracePeriodLt0 := validConfig.DeepCopy()
	nodeMonitorGracePeriodLt0.NodeMonitorGracePeriod = metav1.Duration{Duration: -1 * time.Second}

	nodeStartupGracePeriodLt0 := validConfig.DeepCopy()
	nodeStartupGracePeriodLt0.NodeStartupGracePeriod = metav1.Duration{Duration: -1 * time.Second}

	podEvictionTimeoutLt0 := validConfig.DeepCopy()
	podEvictionTimeoutLt0.PodEvictionTimeout = metav1.Duration{Duration: -1 * time.Second}

	scenarios := map[string]struct {
		expectedToFail bool
		config         *config.NodeLifecycleControllerConfiguration
	}{
		"good": {
			expectedToFail: false,
			config:         validConfig,
		},
		"node-monitor-grace-period-lt-0": {
			expectedToFail: true,
			config:         nodeMonitorGracePeriodLt0,
		},
		"node-start-grace-period-lt-0": {
			expectedToFail: true,
			config:         nodeStartupGracePeriodLt0,
		},
		"pod-eviction-timeout-lt-0": {
			expectedToFail: true,
			config:         podEvictionTimeoutLt0,
		},
	}

	for name, scenario := range scenarios {
		errs := ValidateNodeLifecycleControllerConfiguration(scenario.config, newPath.Child("nodeLifecycleController"))
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}

	}
}

func TestValidatePersistentVolumeBinderControllerConfiguration(t *testing.T) {
	newPath := field.NewPath("KubeControllerManagerConfiguration")
	validConfig := &config.PersistentVolumeBinderControllerConfiguration{
		PVClaimBinderSyncPeriod: metav1.Duration{Duration: 30 * time.Second},
		VolumeConfiguration: config.VolumeConfiguration{
			EnableDynamicProvisioning:  false,
			EnableHostPathProvisioning: true,
			FlexVolumePluginDir:        "/flex-volume-plugin",
			PersistentVolumeRecyclerConfiguration: config.PersistentVolumeRecyclerConfiguration{
				MaximumRetry:             3,
				MinimumTimeoutNFS:        200,
				IncrementTimeoutNFS:      45,
				MinimumTimeoutHostPath:   45,
				IncrementTimeoutHostPath: 45,
			},
		},
	}

	pvClaimBinderSyncPeriodLt0 := validConfig.DeepCopy()
	pvClaimBinderSyncPeriodLt0.PVClaimBinderSyncPeriod = metav1.Duration{Duration: -30 * time.Second}

	scenarios := map[string]struct {
		expectedToFail bool
		config         *config.PersistentVolumeBinderControllerConfiguration
	}{
		"good": {
			expectedToFail: false,
			config:         validConfig,
		},
		"pv-claim-binder-sync-period-lt-0": {
			expectedToFail: true,
			config:         pvClaimBinderSyncPeriodLt0,
		},
	}

	for name, scenario := range scenarios {
		errs := ValidatePersistentVolumeBinderControllerConfiguration(scenario.config, newPath.Child("persistentVolumeBinderController"))
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}

	}
}
