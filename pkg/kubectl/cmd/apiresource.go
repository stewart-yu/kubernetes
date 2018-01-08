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

package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
	"fmt"
	"io"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

var (
	apiresourcesExample = templates.Examples(i18n.T(`
		# Print the supported API resource
		kubectl apiresource`))
)

// apiresourcesOptions: describe the options available to users of the "kubectl
// apiresources" command.
type ApiResourcesOptions struct {
	out io.Writer

	namespaced bool
	apiGroup   string
	output     string
}

func NewCmdApiResources(f cmdutil.Factory, out io.Writer) *cobra.Command {

	options := &ApiResourcesOptions{
		out: out,
	}

	cmd := &cobra.Command{
		Use:     "apiresources",
		Short:   i18n.T("List all resources with different types"),
		Long:    "List all resources with different types in the cluster",
		Example: apiresourcesExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(cmd))
			cmdutil.CheckErr(options.ValidateArgs(cmd, args))
			cmdutil.CheckErr(options.RunApiResources(cmd, f))
		},
	}

	cmd.Flags().BoolP("namespaced", "", false, "Namespaced indicates if a api resource is namespaced or not.")
	//cmd.Flags().StringP("output", "o", "", "Output mode. Use \"-o wide\" for wide output.")
	cmd.Flags().StringP("apigroup", "", "", "The API group to use when talking to the server.")
	cmd.Flags().StringP("output", "o", "", "One of 'yaml' or 'json'.")
	return cmd
}

func (o *ApiResourcesOptions) Complete(cmd *cobra.Command) error {
	o.namespaced = cmdutil.GetFlagBool(cmd, "namespaced")
	o.apiGroup = cmdutil.GetFlagString(cmd, "apigroup")
	o.output = cmdutil.GetFlagString(cmd, "output")
	return nil
}

func (o *ApiResourcesOptions) ValidateArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("not expect too many arguments")
	}

	outputMode := cmdutil.GetFlagString(cmd, "output")
	if o.output != "" && o.output != "wide" && o.output != "yaml" && o.output != "json" {
		return fmt.Errorf("unexpected -o output mode: %v. --output should be one of 'yaml'|'wide'|'json'", outputMode)
	}

	return nil
}

func (o *ApiResourcesOptions) RunApiResources(cmd *cobra.Command, f cmdutil.Factory) error {

	return nil
}
