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

import(
	"io"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
	"github.com/spf13/cobra"

	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"fmt"
)

// ApiResourceOptions contains the input to the get command.
type ApiResourcesOptions struct {
	Out io.Writer

	Namespaced     bool
	ApiResource   string
	output     string
}

var (
	apiResourceLong = templates.LongDesc(`
		Print the supported API resource`)

	apiResourcesExample = templates.Examples(i18n.T(`
		# List all apisources with table.
		kubectl apiresources

		# List all apisources with more information (such as Verbs set).
		kubectl apiresources -o wide

		# List a single apiresource with specified NAME.
		kubectl apiresources apiresource ***

		# List a single apiresource in JSON/YAML output format.
		kubectl apiresources -o json apiresource ***

		# List all *** and *** together.
		kubectl apiresources apiresource ***,***

		# List all apiresources marked namespace.
		kubectl apiresources namespaced=true`))
)

// NewCmdApiResources creates a command object for the generic "apiresource" action, which
// print all apiresource in cluster.
func NewCmdApiResources(f cmdutil.Factory, out io.Writer) *cobra.Command {
	options := &ApiResourcesOptions{
		Out:    out,
	}

	cmd := &cobra.Command{
		Use:     "apiresource [(-o|--output=)json|yaml|wide] ([APIRESOURCETYPE NAME] | [NAMESPACED TRUE|FALSE] ...) [flags]",
		Short:   i18n.T("List all resources with different types"),
		Long:    apiResourceLong + "\n\n" + cmdutil.ValidResourceTypeList(f),
		Example: apiResourcesExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.ValidateArgs(cmd, args))
			cmdutil.CheckErr(options.RunApiResources(f, cmd, args))
		},
	}

	cmd.Flags().BoolVar(&options.Namespaced, "namespaced", options.Namespaced, "If present, list the resource type for the requested object(s).")
	cmd.Flags().StringVarP(&options.ApiResource, "apiresource", "ar", options.ApiResource, " from a server.")

	return cmd
}

// Validate checks the set of flags provided by the user.
func (options *ApiResourcesOptions) ValidateArgs(cmd *cobra.Command, args []string) error {

	outputMode := cmdutil.GetFlagString(cmd, "output")
	if options.output != "" && options.output != "wide" && options.output != "yaml" && options.output != "json" {
		return fmt.Errorf("unexpected -o output mode: %v. --output should be one of 'yaml'|'wide'|'json'", outputMode)
	}
	return nil
}

func (o *ApiResourcesOptions) RunApiResources(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	return nil
}