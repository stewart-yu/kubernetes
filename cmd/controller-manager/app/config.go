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

package app

import (
	apiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/kubernetes/pkg/apis/componentconfig"
)

type ExtraConfig struct {
	Master     string
	Kubeconfig string
}

// Config is the main context object for the controller manager.
type Config struct {
	ComponentConfig   componentconfig.KubeControllerManagerConfiguration
	SecureServingInfo apiserver.SecureServingInfo
	Extra             ExtraConfig
}

type completedConfig struct {
	ComponentConfig   componentconfig.KubeControllerManagerConfiguration
	SecureServingInfo apiserver.SecureServingInfo
	Extra             *ExtraConfig
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	cc := completedConfig{
		c.ComponentConfig,
		c.SecureServingInfo,
		&c.Extra,
	}

	return CompletedConfig{&cc}
}
