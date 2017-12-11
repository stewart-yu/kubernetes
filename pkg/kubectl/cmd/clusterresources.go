/*
Copyright 2017 The Kubernetes Authors.
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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/printers"
)

var (
	clusterresourcesExample = templates.Examples(i18n.T(`
		# List all resources with different types in cluster
		kubectl clusterresources`))
)

// groupResource contains the APIGroup and APIResource
type groupResource struct {
	APIGroup    string
	APIResource metav1.APIResource
}

// clusterresourcesOptions: describe the options available to users of the "kubectl
// clusterresources" command.
type ClusterResourcesOptions struct {
	out io.Writer

	namespaced bool
	apiGroup   string
	output     string
}

func NewCmdClusterResources(f cmdutil.Factory, out io.Writer) *cobra.Command {

	options := &ClusterResourcesOptions{
		out: out,
	}

	cmd := &cobra.Command{
		Use:     "clusterresources",
		Short:   i18n.T("List all resources with different types"),
		Long:    "List all resources with different types in the cluster",
		Example: clusterresourcesExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(cmd))
			cmdutil.CheckErr(options.Validate(args))
			cmdutil.CheckErr(options.RunClusterResources(cmd, f))
		},
	}

	cmd.Flags().BoolP("namespaced", "", false, "Namespaced indicates if a resource is namespaced or not.")
	//cmd.Flags().StringP("output", "o", "", "Output mode. Use \"-o wide\" for wide output.")
	cmd.Flags().StringP("apigroup", "", "", "The API group to use when talking to the server.")
	cmd.Flags().StringP("output", "o", "", "One of 'yaml' or 'json'.")
	return cmd
}

func (o *ClusterResourcesOptions) GetClusterResource(f cmdutil.Factory) ([]groupResource, error) {
	discoveryClient, err := f.DiscoveryClient()
	if err != nil {
		return nil, err
	}

	// Always request fresh data from the server
	discoveryClient.Invalidate()

	lists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, fmt.Errorf("get available api resources from server failed: %v", err)
	}

	resources := []groupResource{}

	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			glog.V(1).Infof("Unable to parse groupversion %s:%s", list.GroupVersion, err.Error())
			continue
		}
		for _, resource := range list.APIResources {
			if len(resource.Verbs) == 0 {
				continue
			}
			//// filter apiGroup
			//if o.apiGroup != gv.Group {
			//	continue
			//}
			//// filter namespaced
			//if o.namespaced != resource.Namespaced {
			//	continue
			//}
			resources = append(resources, groupResource{
				APIGroup:    gv.Group,
				APIResource: resource,
			})
		}
	}

	return resources, nil
}

func (o *ClusterResourcesOptions) RunClusterResources(cmd *cobra.Command, f cmdutil.Factory) error {
	w := printers.GetNewTabWriter(o.out)
	defer w.Flush()

	resources, err := o.GetClusterResource(f)
	if err != nil {
		glog.V(1).Infof("Get available cluster resources failed: %v", err)
	}

	sort.Stable(sortableGroupResource(resources))

	var resultResources string
	// print output
	switch o.output {
	case "":
		resultResources = "NAME\tSHORTNAMES\tAPIGROUP\tNAMESPACED\tKIND\n"
		for _, r := range resources {
			resultResources = fmt.Sprintf("%s%s\t%s\t%s\t%v\t%s\n",
				resultResources,
				r.APIResource.Name,
				strings.Join(r.APIResource.ShortNames, ","),
				r.APIGroup,
				r.APIResource.Namespaced,
				r.APIResource.Kind)
		}
		// normal
	case "wide":
		resultResources = "NAME\tSHORTNAMES\tAPIGROUP\tNAMESPACED\tKIND\tVERBS\n"
		for _, r := range resources {
			resultResources = fmt.Sprintf("%s%s\t%s\t%s\t%v\t%s\t%v\n",
				resultResources,
				r.APIResource.Name,
				strings.Join(r.APIResource.ShortNames, ","),
				r.APIGroup,
				r.APIResource.Namespaced,
				r.APIResource.Kind,
				r.APIResource.Verbs)
		}
	case "yaml":
		// output with yaml
		for _, r := range resources {
			marshalled, err := yaml.Marshal(&r.APIResource)
			if err != nil {
				return err
			}
			resultResources = fmt.Sprintf("%s%s\n", resultResources, string(marshalled))
		}
	case "json":
		// output with json
		for _, r := range resources {
			marshalled, err := json.MarshalIndent(&r.APIResource, "", "  ")
			if err != nil {
				return err
			}
			resultResources = fmt.Sprintf("%s%s\n", string(marshalled))
		}
	default:
		// There is a bug in the program if we hit this case.
		// However, we follow a policy of never panicking.
		return fmt.Errorf("versionOptions were not validated: --output=%q should have been rejected", o.output)
	}
	if len(resources) > 0{
		fmt.Fprintf(w, "%v", resultResources)
	}else{
		fmt.Fprintln(w, "No resources found.")
	}

	return nil
}

func (o *ClusterResourcesOptions) Complete(cmd *cobra.Command) error {
	o.namespaced = cmdutil.GetFlagBool(cmd, "namespaced")
	o.apiGroup = cmdutil.GetFlagString(cmd, "apigroup")
	o.output = cmdutil.GetFlagString(cmd, "output")
	return nil
}

func (o *ClusterResourcesOptions) Validate(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("not expect too many arguments")
	}

	if o.output != "" && o.output != "wide" && o.output != "yaml" && o.output != "json" {
		return fmt.Errorf("unexpected -o output mode: %v. --output should be one of 'yaml'|'wide'|'json'", )
	}

	return nil
}

type sortableGroupResource []groupResource

func (s sortableGroupResource) Len() int      { return len(s) }
func (s sortableGroupResource) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortableGroupResource) Less(i, j int) bool {
	ret := strings.Compare(s[i].APIGroup, s[j].APIGroup)
	if ret > 0 {
		return false
	} else if ret == 0 {
		return strings.Compare(s[i].APIResource.Name, s[j].APIResource.Name) < 0
	}
	return true
}
